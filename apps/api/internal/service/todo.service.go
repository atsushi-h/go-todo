package service

import (
	"context"
	"errors"

	"go-todo/db/sqlc"

	"github.com/jackc/pgx/v5"
)

var ErrTodoNotFound = errors.New("todo not found")

// バッチ処理の結果
type BatchCompleteResult struct {
	Succeeded []sqlc.Todo
	Failed    []BatchFailedItem
}

type BatchDeleteResult struct {
	Succeeded []int64
	Failed    []BatchFailedItem
}

type BatchFailedItem struct {
	ID    int64
	Error string
}

type TodoService struct {
	repo TodoRepository
}

func NewTodoService(repo TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) GetAllTodos(ctx context.Context, userID int64) ([]sqlc.Todo, error) {
	return s.repo.ListTodosByUser(ctx, userID)
}

func (s *TodoService) GetTodoByID(ctx context.Context, id, userID int64) (*sqlc.Todo, error) {
	todo, err := s.repo.GetTodoByID(ctx, sqlc.GetTodoByIDParams{
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
	todo, err := s.repo.CreateTodo(ctx, sqlc.CreateTodoParams{
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
	todo, err := s.repo.UpdateTodo(ctx, sqlc.UpdateTodoParams{
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
	return s.repo.DeleteTodo(ctx, sqlc.DeleteTodoParams{
		ID:     id,
		UserID: userID,
	})
}

func (s *TodoService) BatchCompleteTodos(ctx context.Context, userID int64, ids []int64) (*BatchCompleteResult, error) {
	result := &BatchCompleteResult{
		Succeeded: []sqlc.Todo{},
		Failed:    []BatchFailedItem{},
	}

	// 存在チェック
	existingTodos, err := s.repo.GetTodosByIDs(ctx, sqlc.GetTodosByIDsParams{
		Ids:    ids,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	// 存在するIDのマップを作成
	existingIDMap := make(map[int64]bool)
	for _, todo := range existingTodos {
		existingIDMap[todo.ID] = true
	}

	// 存在しないIDを失敗として記録
	var validIDs []int64
	for _, id := range ids {
		if existingIDMap[id] {
			validIDs = append(validIDs, id)
		} else {
			result.Failed = append(result.Failed, BatchFailedItem{
				ID:    id,
				Error: "Todo not found",
			})
		}
	}

	// 有効なIDがある場合のみバッチ更新
	if len(validIDs) > 0 {
		completedTodos, err := s.repo.BatchCompleteTodos(ctx, sqlc.BatchCompleteTodosParams{
			Ids:    validIDs,
			UserID: userID,
		})
		if err != nil {
			return nil, err
		}
		result.Succeeded = completedTodos
	}

	return result, nil
}

func (s *TodoService) BatchDeleteTodos(ctx context.Context, userID int64, ids []int64) (*BatchDeleteResult, error) {
	result := &BatchDeleteResult{
		Succeeded: []int64{},
		Failed:    []BatchFailedItem{},
	}

	// 存在チェック
	existingTodos, err := s.repo.GetTodosByIDs(ctx, sqlc.GetTodosByIDsParams{
		Ids:    ids,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	existingIDMap := make(map[int64]bool)
	for _, todo := range existingTodos {
		existingIDMap[todo.ID] = true
	}

	var validIDs []int64
	for _, id := range ids {
		if existingIDMap[id] {
			validIDs = append(validIDs, id)
		} else {
			result.Failed = append(result.Failed, BatchFailedItem{
				ID:    id,
				Error: "Todo not found",
			})
		}
	}

	// バッチ削除実行
	if len(validIDs) > 0 {
		err := s.repo.BatchDeleteTodos(ctx, sqlc.BatchDeleteTodosParams{
			Ids:    validIDs,
			UserID: userID,
		})
		if err != nil {
			return nil, err
		}
		result.Succeeded = validIDs
	}

	return result, nil
}
