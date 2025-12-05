package model

import "time"

// Todo構造体
type Todo struct {
	ID          uint      `json:"id" bun:",pk,autoincrement"`
	Title       string    `json:"title" bun:",notnull"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed" bun:",default:false"`
	CreatedAt   time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
	UserID      uint      `json:"user_id" bun:",notnull"`
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
