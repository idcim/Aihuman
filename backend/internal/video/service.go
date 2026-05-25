package video

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"strings"
	"time"
)

var (
	ErrInvalidDuration = errors.New("duration_seconds must be between 15 and 180")
	ErrTaskNotFound    = errors.New("task not found")
)

type Service struct {
	store  Store
	queue  Queue
	logger *slog.Logger
}

func NewService(store Store, logger *slog.Logger) *Service {
	return &Service{store: store, logger: logger}
}

func (s *Service) UseQueue(queue Queue) {
	s.queue = queue
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (Task, error) {
	duration := req.DurationSeconds
	if duration == 0 {
		duration = 60
	}
	if duration < 15 || duration > 180 {
		return Task{}, ErrInvalidDuration
	}

	now := time.Now().UTC()
	task := Task{
		ID:              newID(),
		ProductName:     strings.TrimSpace(req.ProductName),
		Script:          strings.TrimSpace(req.Script),
		AvatarImageURL:  strings.TrimSpace(req.AvatarImageURL),
		DurationSeconds: duration,
		Status:          TaskStatusQueued,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if task.ProductName == "" || task.Script == "" || task.AvatarImageURL == "" {
		return Task{}, errors.New("product_name, script, and avatar_image_url are required")
	}

	if err := s.store.Save(ctx, task); err != nil {
		return Task{}, err
	}
	if s.queue != nil {
		if err := s.queue.Enqueue(ctx, task.ID); err != nil {
			return Task{}, err
		}
	}
	s.logger.Info("video task queued", "task_id", task.ID)
	return task, nil
}

func (s *Service) Get(ctx context.Context, id string) (Task, error) {
	return s.store.Get(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Task, error) {
	return s.store.List(ctx)
}

func (s *Service) MarkRunning(ctx context.Context, id string) (Task, error) {
	return s.store.UpdateStatus(ctx, id, TaskStatusRunning, RenderResult{}, "")
}

func (s *Service) MarkSucceeded(ctx context.Context, id string, result RenderResult) (Task, error) {
	return s.store.UpdateStatus(ctx, id, TaskStatusSucceeded, result, "")
}

func (s *Service) MarkFailed(ctx context.Context, id string, errorMessage string) (Task, error) {
	return s.store.UpdateStatus(ctx, id, TaskStatusFailed, RenderResult{}, errorMessage)
}

func newID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(bytes[:])
}
