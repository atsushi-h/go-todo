package repository

import (
	"errors"
	"sync"
	"time"

	"go-todo/model"
)

var (
	ErrTodoNotFound = errors.New("todo not found")
)

// Todoのデータアクセスを管理
type TodoRepository struct {
	mu      sync.RWMutex
	todos   map[int]*model.Todo
	nextID  int
}

// 新しいTodoRepositoryを作成
func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		todos:  make(map[int]*model.Todo),
		nextID: 1,
	}
}

// 全てのTodoを取得
func (r *TodoRepository) GetAll() []*model.Todo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos := make([]*model.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}
	return todos
}

// 指定されたIDのTodoを取得
func (r *TodoRepository) GetByID(id int) (*model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, exists := r.todos[id]
	if !exists {
		return nil, ErrTodoNotFound
	}
	return todo, nil
}

// 新しいTodoを作成
func (r *TodoRepository) Create(title, description string) *model.Todo {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	todo := &model.Todo{
		ID:          r.nextID,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	r.todos[r.nextID] = todo
	r.nextID++

	return todo
}

// 既存のTodoを更新
func (r *TodoRepository) Update(id int, title *string, description *string, completed *bool) (*model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo, exists := r.todos[id]
	if !exists {
		return nil, ErrTodoNotFound
	}

	if title != nil {
		todo.Title = *title
	}
	if description != nil {
		todo.Description = *description
	}
	if completed != nil {
		todo.Completed = *completed
	}
	todo.UpdatedAt = time.Now()

	return todo, nil
}

// 指定されたIDのTodoを削除
func (r *TodoRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[id]; !exists {
		return ErrTodoNotFound
	}

	delete(r.todos, id)
	return nil
}
