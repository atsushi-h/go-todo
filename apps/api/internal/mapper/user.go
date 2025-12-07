package mapper

import (
	"go-todo/db/sqlc"
)

type UserResponse struct {
	ID        int64   `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

func UserToResponse(u *sqlc.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		AvatarURL: u.AvatarUrl,
	}
}
