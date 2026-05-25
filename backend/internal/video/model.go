package video

import "time"

type CreateRequest struct {
	ProductName     string `json:"product_name" binding:"required"`
	Script          string `json:"script" binding:"required"`
	AvatarImageURL  string `json:"avatar_image_url" binding:"required"`
	DurationSeconds int    `json:"duration_seconds"`
}

type TaskStatus string

const (
	TaskStatusQueued TaskStatus = "queued"
	TaskStatusFailed TaskStatus = "failed"
)

type Task struct {
	ID              string     `json:"id"`
	ProductName     string     `json:"product_name"`
	Script          string     `json:"script"`
	AvatarImageURL  string     `json:"avatar_image_url"`
	DurationSeconds int        `json:"duration_seconds"`
	Status          TaskStatus `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
