package model

import "time"

// Todo構造体
type Todo struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Title       string    `json:"title" gorm:"not null"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed" gorm:"default:false"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
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
