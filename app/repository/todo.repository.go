package repository

import (
	"errors"

	"go-todo/model"
	"gorm.io/gorm"
)

type TodoRepository interface {
    GetAll() ([]*model.Todo, error)
    GetByID(id uint) (*model.Todo, error)
    Create(title, description string) (*model.Todo, error)
    Update(id uint, title *string, description *string, completed *bool) (*model.Todo, error)
    Delete(id uint) error
}

type todoRepository struct {
    db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
    return &todoRepository{db: db}
}

var (
	ErrTodoNotFound = errors.New("todo not found")
)

// 全てのTodoを取得
func (r *todoRepository) GetAll() ([]*model.Todo, error) {
	var todos []*model.Todo

    if err := r.db.Order("created_at DESC").Find(&todos).Error; err != nil {
        return nil, err
    }

    return todos, nil
}

// 指定されたIDのTodoを取得
func (r *todoRepository) GetByID(id uint) (*model.Todo, error) {
    var todo model.Todo

    if err := r.db.First(&todo, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrTodoNotFound
        }
        return nil, err
    }

    return &todo, nil
}

// 新しいTodoを作成
func (r *todoRepository) Create(title, description string) (*model.Todo, error) {
    todo := &model.Todo{
        Title:       title,
        Description: description,
        Completed:   false,
    }
    
    if err := r.db.Create(todo).Error; err != nil {
        return nil, err
    }

    return todo, nil
}

// 既存のTodoを更新
func (r *todoRepository) Update(id uint, title *string, description *string, completed *bool) (*model.Todo, error) {
    var todo model.Todo

    // レコードの存在確認
    if err := r.db.First(&todo, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrTodoNotFound
        }
        return nil, err
    }

    // 更新データの準備
    updates := make(map[string]interface{})
    if title != nil {
        updates["title"] = *title
    }
    if description != nil {
        updates["description"] = *description
    }
    if completed != nil {
        updates["completed"] = *completed
    }

    // 更新実行
    if err := r.db.Model(&todo).Updates(updates).Error; err != nil {
        return nil, err
    }

    return &todo, nil
}

// 指定されたIDのTodoを削除
func (r *todoRepository) Delete(id uint) error {
    result := r.db.Delete(&model.Todo{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return ErrTodoNotFound
    }
    return nil
}
