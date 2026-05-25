package video

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Worker struct {
	service        *Service
	queue          Queue
	videoEngineURL string
	httpClient     *http.Client
	logger         *slog.Logger
}

func NewWorker(service *Service, queue Queue, videoEngineURL string, logger *slog.Logger) *Worker {
	return &Worker{
		service:        service,
		queue:          queue,
		videoEngineURL: videoEngineURL,
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		logger:         logger,
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.logger.Info("video worker started")
	go w.loop(ctx)
}

func (w *Worker) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			w.logger.Info("video worker stopped")
			return
		default:
		}

		queued, err := w.queue.Consume(ctx)
		if errors.Is(err, ErrQueueEmpty) {
			continue
		}
		if err != nil {
			w.logger.Warn("failed to consume video task", "error", err)
			time.Sleep(time.Second)
			continue
		}

		w.process(ctx, queued)
	}
}

func (w *Worker) process(ctx context.Context, queued QueuedTask) {
	task, err := w.service.MarkRunning(ctx, queued.TaskID)
	if err != nil {
		w.logger.Warn("failed to mark task running", "task_id", queued.TaskID, "error", err)
		return
	}

	result, err := w.renderVideo(ctx, task)
	if err != nil {
		_, _ = w.service.MarkFailed(ctx, task.ID, err.Error())
		w.logger.Warn("video task failed", "task_id", task.ID, "error", err)
		_ = w.queue.Ack(ctx, queued.MessageID)
		return
	}

	if _, err := w.service.MarkSucceeded(ctx, task.ID, result); err != nil {
		w.logger.Warn("failed to mark task succeeded", "task_id", task.ID, "error", err)
		return
	}
	if err := w.queue.Ack(ctx, queued.MessageID); err != nil {
		w.logger.Warn("failed to ack video task", "task_id", task.ID, "error", err)
		return
	}
	w.logger.Info("video task completed", "task_id", task.ID)
}

func (w *Worker) renderVideo(ctx context.Context, task Task) (RenderResult, error) {
	payload := renderRequest{
		TaskID:          task.ID,
		ProductName:     task.ProductName,
		Script:          task.Script,
		AvatarImagePath: task.AvatarImageURL,
		VoiceoverPath:   fmt.Sprintf("/tmp/%s.wav", task.ID),
		OutputPath:      fmt.Sprintf("/output/%s.mp4", task.ID),
		DurationSeconds: task.DurationSeconds,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return RenderResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.videoEngineURL+"/v1/render", bytes.NewReader(body))
	if err != nil {
		return RenderResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return RenderResult{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return RenderResult{}, fmt.Errorf("video engine returned HTTP %d", resp.StatusCode)
	}

	var response renderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return RenderResult{}, err
	}
	return RenderResult{
		Command:        response.Command,
		OutputURL:      response.ResultURL,
		ObjectKey:      response.ObjectKey,
		PreviewSeconds: response.PreviewDurationSeconds,
	}, nil
}

type renderRequest struct {
	TaskID          string `json:"task_id"`
	ProductName     string `json:"product_name"`
	Script          string `json:"script"`
	AvatarImagePath string `json:"avatar_image_path"`
	VoiceoverPath   string `json:"voiceover_path"`
	OutputPath      string `json:"output_path"`
	DurationSeconds int    `json:"duration_seconds"`
}

type renderResponse struct {
	Command                []string `json:"command"`
	OutputPath             string   `json:"output_path"`
	ObjectKey              string   `json:"object_key"`
	ResultURL              string   `json:"result_url"`
	PreviewDurationSeconds int      `json:"preview_duration_seconds"`
}
