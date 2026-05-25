package video

import (
	"context"
	"sort"
	"sync"
	"time"
)

type Store interface {
	Save(ctx context.Context, task Task) error
	Get(ctx context.Context, id string) (Task, error)
	List(ctx context.Context) ([]Task, error)
	UpdateStatus(ctx context.Context, id string, status TaskStatus, result RenderResult, errorMessage string) (Task, error)
}

type MemoryStore struct {
	mu    sync.RWMutex
	tasks map[string]Task
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{tasks: make(map[string]Task)}
}

func (s *MemoryStore) Save(ctx context.Context, task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
	return nil
}

func (s *MemoryStore) Get(ctx context.Context, id string) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (s *MemoryStore) List(ctx context.Context) ([]Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})
	return tasks, nil
}

func (s *MemoryStore) UpdateStatus(ctx context.Context, id string, status TaskStatus, result RenderResult, errorMessage string) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}
	task.Status = status
	task.RenderCommand = result.Command
	task.OutputURL = result.OutputURL
	task.ObjectKey = result.ObjectKey
	task.PreviewSeconds = result.PreviewSeconds
	task.ErrorMessage = errorMessage
	task.UpdatedAt = time.Now().UTC()
	s.tasks[id] = task
	return task, nil
}
