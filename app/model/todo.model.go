package model

import "time"

// Todo構造体
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Todo作成リクエストの構造体
type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Todo更新リクエストの構造体
type UpdateTodoRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}
