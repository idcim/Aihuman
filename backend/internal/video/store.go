package video

import (
	"context"
	"sync"
)

type Store interface {
	Save(ctx context.Context, task Task) error
	Get(ctx context.Context, id string) (Task, error)
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
