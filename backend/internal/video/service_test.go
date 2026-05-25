package video

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func TestServiceCreateDefaultsDurationAndQueuesTask(t *testing.T) {
	service := NewService(NewMemoryStore(), slog.New(slog.NewTextHandler(io.Discard, nil)))

	task, err := service.Create(context.Background(), CreateRequest{
		ProductName:    "Demo product",
		Script:         "Demo script",
		AvatarImageURL: "https://example.com/avatar.png",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if task.DurationSeconds != 60 {
		t.Fatalf("DurationSeconds = %d, want 60", task.DurationSeconds)
	}
	if task.Status != TaskStatusQueued {
		t.Fatalf("Status = %s, want %s", task.Status, TaskStatusQueued)
	}
	if task.ID == "" {
		t.Fatal("ID is empty")
	}
}

func TestServiceCreateRejectsInvalidDuration(t *testing.T) {
	service := NewService(NewMemoryStore(), slog.New(slog.NewTextHandler(io.Discard, nil)))

	_, err := service.Create(context.Background(), CreateRequest{
		ProductName:     "Demo product",
		Script:          "Demo script",
		AvatarImageURL:  "https://example.com/avatar.png",
		DurationSeconds: 10,
	})
	if err != ErrInvalidDuration {
		t.Fatalf("Create error = %v, want %v", err, ErrInvalidDuration)
	}
}

func TestServiceGetReturnsStoredTask(t *testing.T) {
	service := NewService(NewMemoryStore(), slog.New(slog.NewTextHandler(io.Discard, nil)))

	created, err := service.Create(context.Background(), CreateRequest{
		ProductName:    "Demo product",
		Script:         "Demo script",
		AvatarImageURL: "https://example.com/avatar.png",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	got, err := service.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("Get ID = %s, want %s", got.ID, created.ID)
	}
}
