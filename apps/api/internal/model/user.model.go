package model

import "time"

type User struct {
	ID         uint      `json:"id" bun:",pk,autoincrement"`
	Email      string    `json:"email" bun:",unique,notnull"`
	Name       string    `json:"name"`
	AvatarURL  string    `json:"avatar_url"`
	Provider   string    `json:"provider" bun:",notnull"`
	ProviderID string    `json:"-" bun:",notnull"`
	CreatedAt  time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt  time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
}
