package service

import (
	"context"

	"go-todo/db/sqlc"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int64) (sqlc.User, error)
	GetUserByProviderID(ctx context.Context, arg sqlc.GetUserByProviderIDParams) (sqlc.User, error)
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
	DeleteUser(ctx context.Context, id int64) error
	DeleteTodosByUserID(ctx context.Context, userID int64) error
}

// sqlc.Querier が UserRepository を満たすことを保証
var _ UserRepository = (sqlc.Querier)(nil)
