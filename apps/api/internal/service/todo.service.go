package service

import (
	"context"
	"errors"

	"go-todo/db/sqlc"

	"github.com/jackc/pgx/v5"
)

var ErrTodoNotFound = errors.New("todo not found")

type TodoService struct {
	querier sqlc.Querier
}

func NewTodoService(querier sqlc.Querier) *TodoService {
	return &TodoService{querier: querier}
}

func (s *TodoService) GetAllTodos(ctx context.Context, userID int64) ([]sqlc.Todo, error) {
	return s.querier.ListTodosByUser(ctx, userID)
}

func (s *TodoService) GetTodoByID(ctx context.Context, id, userID int64) (*sqlc.Todo, error) {
	todo, err := s.querier.GetTodoByID(ctx, sqlc.GetTodoByIDParams{
		ID:     id,
		UserID: userID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTodoNotFound
	}
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (s *TodoService) CreateTodo(ctx context.Context, userID int64, title string, description *string) (*sqlc.Todo, error) {
	todo, err := s.querier.CreateTodo(ctx, sqlc.CreateTodoParams{
		UserID:      userID,
		Title:       title,
		Description: description,
	})
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, id, userID int64, title, description *string, completed *bool) (*sqlc.Todo, error) {
	todo, err := s.querier.UpdateTodo(ctx, sqlc.UpdateTodoParams{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: description,
		Completed:   completed,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTodoNotFound
	}
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, id, userID int64) error {
	return s.querier.DeleteTodo(ctx, sqlc.DeleteTodoParams{
		ID:     id,
		UserID: userID,
	})
}
