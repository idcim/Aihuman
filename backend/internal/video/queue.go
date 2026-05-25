package video

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Queue interface {
	Enqueue(ctx context.Context, taskID string) error
	Consume(ctx context.Context) (QueuedTask, error)
	Ack(ctx context.Context, messageID string) error
	Close() error
}

type QueuedTask struct {
	MessageID string
	TaskID    string
}

var ErrQueueEmpty = errors.New("queue empty")

type RedisQueue struct {
	client *redis.Client
	stream string
	group  string
	worker string
}

func NewRedisQueue(ctx context.Context, addr string, stream string, logger *slog.Logger) (*RedisQueue, error) {
	client := redis.NewClient(&redis.Options{Addr: addr})
	if err := waitForRedis(ctx, client, 30*time.Second); err != nil {
		_ = client.Close()
		return nil, err
	}

	queue := &RedisQueue{
		client: client,
		stream: stream,
		group:  "video-workers",
		worker: "backend-worker",
	}

	err := client.XGroupCreateMkStream(ctx, stream, queue.group, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		_ = client.Close()
		return nil, err
	}

	logger.Info("using redis stream queue", "stream", stream, "group", queue.group)
	return queue, nil
}

func (q *RedisQueue) Enqueue(ctx context.Context, taskID string) error {
	return q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: q.stream,
		Values: map[string]interface{}{"task_id": taskID},
	}).Err()
}

func (q *RedisQueue) Consume(ctx context.Context) (QueuedTask, error) {
	streams, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    q.group,
		Consumer: q.worker,
		Streams:  []string{q.stream, ">"},
		Count:    1,
		Block:    5 * time.Second,
	}).Result()
	if errors.Is(err, redis.Nil) {
		return QueuedTask{}, ErrQueueEmpty
	}
	if err != nil {
		return QueuedTask{}, err
	}
	if len(streams) == 0 || len(streams[0].Messages) == 0 {
		return QueuedTask{}, ErrQueueEmpty
	}

	message := streams[0].Messages[0]
	taskID, _ := message.Values["task_id"].(string)
	if taskID == "" {
		return QueuedTask{}, errors.New("queued message missing task_id")
	}

	return QueuedTask{MessageID: message.ID, TaskID: taskID}, nil
}

func (q *RedisQueue) Ack(ctx context.Context, messageID string) error {
	return q.client.XAck(ctx, q.stream, q.group, messageID).Err()
}

func (q *RedisQueue) Close() error {
	return q.client.Close()
}

func waitForRedis(ctx context.Context, client *redis.Client, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		if err := client.Ping(ctx).Err(); err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
