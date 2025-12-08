package service

import (
	"context"

	"go-todo/db/sqlc"
)

type TodoRepository interface {
	GetTodoByID(ctx context.Context, arg sqlc.GetTodoByIDParams) (sqlc.Todo, error)
	ListTodosByUser(ctx context.Context, userID int64) ([]sqlc.Todo, error)
	CreateTodo(ctx context.Context, arg sqlc.CreateTodoParams) (sqlc.Todo, error)
	UpdateTodo(ctx context.Context, arg sqlc.UpdateTodoParams) (sqlc.Todo, error)
	DeleteTodo(ctx context.Context, arg sqlc.DeleteTodoParams) error
}

// sqlc.Querier が TodoRepository を満たすことを保証
var _ TodoRepository = (sqlc.Querier)(nil)
