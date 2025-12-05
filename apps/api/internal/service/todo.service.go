package service

import (
	"context"
	"errors"

	"go-todo/internal/model"
	"go-todo/internal/repository"
)

// Service層のエラー定義
var (
	ErrTodoNotFound = errors.New("todo not found")
)

type TodoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

// 全てのTodoを取得
func (s *TodoService) GetAllTodos(ctx context.Context, userID uint) ([]*model.Todo, error) {
	todos, err := s.repo.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

// 指定されたIDのTodoを取得
func (s *TodoService) GetTodoByID(ctx context.Context, id uint, userID uint) (*model.Todo, error) {
	todo, err := s.repo.GetByID(ctx, id, userID)
	if err == repository.ErrTodoNotFound {
		return nil, ErrTodoNotFound
	}
	return todo, err
}

// 新しいTodoを作成
func (s *TodoService) CreateTodo(ctx context.Context, req model.CreateTodoRequest, userID uint) (*model.Todo, error) {
	todo, err := s.repo.Create(ctx, req.Title, req.Description, userID)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// 既存のTodoを更新
func (s *TodoService) UpdateTodo(ctx context.Context, id uint, userID uint, req model.UpdateTodoRequest) (*model.Todo, error) {
	todo, err := s.repo.Update(ctx, id, userID, req.Title, req.Description, req.Completed)
	if err == repository.ErrTodoNotFound {
		return nil, ErrTodoNotFound
	}
	return todo, err
}

// 指定されたIDのTodoを削除
func (s *TodoService) DeleteTodo(ctx context.Context, id uint, userID uint) error {
	err := s.repo.Delete(ctx, id, userID)
	if err == repository.ErrTodoNotFound {
		return ErrTodoNotFound
	}
	return err
}
