package service

import (
	"errors"

	"go-todo/model"
	"go-todo/repository"
)

// Service層のエラー定義
var (
	ErrTodoNotFound = errors.New("todo not found")
)

type TodoService struct {
	repo *repository.TodoRepository
}

func NewTodoService(repo *repository.TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

// 全てのTodoを取得
func (s *TodoService) GetAllTodos() ([]*model.Todo, error) {
	return s.repo.GetAll(), nil
}

// 指定されたIDのTodoを取得
func (s *TodoService) GetTodoByID(id int) (*model.Todo, error) {
	todo, err := s.repo.GetByID(id)
	if err == repository.ErrTodoNotFound {
		return nil, ErrTodoNotFound
	}
	return todo, err
}

// 新しいTodoを作成
func (s *TodoService) CreateTodo(req model.CreateTodoRequest) (*model.Todo, error) {
	return s.repo.Create(req.Title, req.Description), nil
}

// 既存のTodoを更新
func (s *TodoService) UpdateTodo(id int, req model.UpdateTodoRequest) (*model.Todo, error) {
	todo, err := s.repo.Update(id, req.Title, req.Description, req.Completed)
	if err == repository.ErrTodoNotFound {
		return nil, ErrTodoNotFound
	}
	return todo, err
}

// 指定されたIDのTodoを削除
func (s *TodoService) DeleteTodo(id int) error {
	err := s.repo.Delete(id)
	if err == repository.ErrTodoNotFound {
		return ErrTodoNotFound
	}
	return err
}
