package repository

import (
	"context"
	"database/sql"
	"errors"

	"go-todo/internal/model"

	"github.com/uptrace/bun"
)

type TodoRepository interface {
	GetAll(ctx context.Context, userID uint) ([]*model.Todo, error)
	GetByID(ctx context.Context, id uint, userID uint) (*model.Todo, error)
	Create(ctx context.Context, title, description string, userID uint) (*model.Todo, error)
	Update(ctx context.Context, id uint, userID uint, title *string, description *string, completed *bool) (*model.Todo, error)
	Delete(ctx context.Context, id uint, userID uint) error
}

type todoRepository struct {
	db *bun.DB
}

func NewTodoRepository(db *bun.DB) TodoRepository {
	return &todoRepository{db: db}
}

var (
	ErrTodoNotFound = errors.New("todo not found")
)

// 全てのTodoを取得
func (r *todoRepository) GetAll(ctx context.Context, userID uint) ([]*model.Todo, error) {
	var todos []*model.Todo

	if err := r.db.NewSelect().
		Model(&todos).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Scan(ctx); err != nil {
		return nil, err
	}

	return todos, nil
}

// 指定されたIDのTodoを取得
func (r *todoRepository) GetByID(ctx context.Context, id uint, userID uint) (*model.Todo, error) {
	todo := new(model.Todo)

	if err := r.db.NewSelect().
		Model(todo).
		Where("id = ?", id).
		Where("user_id = ?", userID).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}

	return todo, nil
}

// 新しいTodoを作成
func (r *todoRepository) Create(ctx context.Context, title, description string, userID uint) (*model.Todo, error) {
	todo := &model.Todo{
		Title:       title,
		Description: description,
		Completed:   false,
		UserID:      userID,
	}

	if _, err := r.db.NewInsert().
		Model(todo).
		Exec(ctx); err != nil {
		return nil, err
	}

	return todo, nil
}

// 既存のTodoを更新
func (r *todoRepository) Update(ctx context.Context, id uint, userID uint, title *string, description *string, completed *bool) (*model.Todo, error) {
	// レコードの存在確認
	todo, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// 更新データの準備
	if title != nil {
		todo.Title = *title
	}
	if description != nil {
		todo.Description = *description
	}
	if completed != nil {
		todo.Completed = *completed
	}

	// 更新実行
	if _, err := r.db.NewUpdate().
		Model(todo).
		WherePK().
		Exec(ctx); err != nil {
		return nil, err
	}

	return todo, nil
}

// 指定されたIDのTodoを削除
func (r *todoRepository) Delete(ctx context.Context, id uint, userID uint) error {
	result, err := r.db.NewDelete().
		Model((*model.Todo)(nil)).
		Where("id = ?", id).
		Where("user_id = ?", userID).
		Exec(ctx)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTodoNotFound
	}

	return nil
}
