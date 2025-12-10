package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-todo/db/sqlc"
	"go-todo/internal/service/mocks"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTodoService_GetTodoByID(t *testing.T) {
	t.Run("正常系: Todoを取得できる", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)

		ctx := context.Background()
		todoID := int64(1)
		userID := int64(1)
		now := time.Now()

		expectedTodo := sqlc.Todo{
			ID:          todoID,
			UserID:      userID,
			Title:       "Test Todo",
			Description: ptrString("Test Description"),
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mockRepo.EXPECT().
			GetTodoByID(ctx, sqlc.GetTodoByIDParams{
				ID:     todoID,
				UserID: userID,
			}).
			Return(expectedTodo, nil)

		result, err := svc.GetTodoByID(ctx, todoID, userID)

		require.NoError(t, err)
		assert.Equal(t, expectedTodo.ID, result.ID)
		assert.Equal(t, expectedTodo.Title, result.Title)
		assert.Equal(t, expectedTodo.Description, result.Description)
	})

	t.Run("異常系: ErrTodoNotFoundを返す", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)

		ctx := context.Background()
		todoID := int64(999)
		userID := int64(1)

		mockRepo.EXPECT().
			GetTodoByID(ctx, sqlc.GetTodoByIDParams{
				ID:     todoID,
				UserID: userID,
			}).
			Return(sqlc.Todo{}, pgx.ErrNoRows)

		result, err := svc.GetTodoByID(ctx, todoID, userID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrTodoNotFound)
	})

	t.Run("異常系: その他のエラーをそのまま返す", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)

		ctx := context.Background()
		todoID := int64(1)
		userID := int64(1)
		dbErr := errors.New("database error")

		mockRepo.EXPECT().
			GetTodoByID(ctx, sqlc.GetTodoByIDParams{
				ID:     todoID,
				UserID: userID,
			}).
			Return(sqlc.Todo{}, dbErr)

		result, err := svc.GetTodoByID(ctx, todoID, userID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestTodoService_GetAllTodos(t *testing.T) {
	t.Run("正常系: Todo一覧を取得できる", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)

		ctx := context.Background()
		userID := int64(1)
		now := time.Now()

		expectedTodos := []sqlc.Todo{
			{
				ID:        1,
				UserID:    userID,
				Title:     "Todo 1",
				Completed: false,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        2,
				UserID:    userID,
				Title:     "Todo 2",
				Completed: true,
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		mockRepo.EXPECT().
			ListTodosByUser(ctx, userID).
			Return(expectedTodos, nil)

		result, err := svc.GetAllTodos(ctx, userID)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedTodos[0].Title, result[0].Title)
		assert.Equal(t, expectedTodos[1].Title, result[1].Title)
	})

	t.Run("正常系: 空の一覧を返す", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)

		ctx := context.Background()
		userID := int64(1)

		mockRepo.EXPECT().
			ListTodosByUser(ctx, userID).
			Return([]sqlc.Todo{}, nil)

		result, err := svc.GetAllTodos(ctx, userID)

		require.NoError(t, err)
		assert.Len(t, result, 0)
	})
}

func TestTodoService_CreateTodo(t *testing.T) {
	t.Run("正常系: Todoを作成できる", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)

		ctx := context.Background()
		userID := int64(1)
		title := "New Todo"
		description := ptrString("New Description")
		now := time.Now()

		expectedTodo := sqlc.Todo{
			ID:          1,
			UserID:      userID,
			Title:       title,
			Description: description,
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mockRepo.EXPECT().
			CreateTodo(ctx, sqlc.CreateTodoParams{
				UserID:      userID,
				Title:       title,
				Description: description,
			}).
			Return(expectedTodo, nil)

		result, err := svc.CreateTodo(ctx, userID, title, description)

		require.NoError(t, err)
		assert.Equal(t, expectedTodo.ID, result.ID)
		assert.Equal(t, expectedTodo.Title, result.Title)
		assert.Equal(t, expectedTodo.Description, result.Description)
	})

	t.Run("正常系: descriptionなしでTodoを作成できる", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)

		ctx := context.Background()
		userID := int64(1)
		title := "New Todo"
		now := time.Now()

		expectedTodo := sqlc.Todo{
			ID:          1,
			UserID:      userID,
			Title:       title,
			Description: nil,
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mockRepo.EXPECT().
			CreateTodo(ctx, sqlc.CreateTodoParams{
				UserID:      userID,
				Title:       title,
				Description: nil,
			}).
			Return(expectedTodo, nil)

		result, err := svc.CreateTodo(ctx, userID, title, nil)

		require.NoError(t, err)
		assert.Equal(t, expectedTodo.ID, result.ID)
		assert.Nil(t, result.Description)
	})
}

func TestTodoService_UpdateTodo(t *testing.T) {
	t.Run("正常系: Todoを更新できる", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)

		ctx := context.Background()
		todoID := int64(1)
		userID := int64(1)
		newTitle := ptrString("Updated Title")
		newDescription := ptrString("Updated Description")
		completed := ptrBool(true)
		now := time.Now()

		expectedTodo := sqlc.Todo{
			ID:          todoID,
			UserID:      userID,
			Title:       *newTitle,
			Description: newDescription,
			Completed:   *completed,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mockRepo.EXPECT().
			UpdateTodo(ctx, sqlc.UpdateTodoParams{
				ID:          todoID,
				UserID:      userID,
				Title:       newTitle,
				Description: newDescription,
				Completed:   completed,
			}).
			Return(expectedTodo, nil)

		result, err := svc.UpdateTodo(ctx, todoID, userID, newTitle, newDescription, completed)

		require.NoError(t, err)
		assert.Equal(t, *newTitle, result.Title)
		assert.Equal(t, *completed, result.Completed)
	})

	t.Run("異常系: ErrTodoNotFoundを返す", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)

		ctx := context.Background()
		todoID := int64(999)
		userID := int64(1)
		newTitle := ptrString("Updated Title")

		mockRepo.EXPECT().
			UpdateTodo(ctx, mock.Anything).
			Return(sqlc.Todo{}, pgx.ErrNoRows)

		result, err := svc.UpdateTodo(ctx, todoID, userID, newTitle, nil, nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrTodoNotFound)
	})
}

func TestTodoService_DeleteTodo(t *testing.T) {
	t.Run("正常系: Todoを削除できる", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)

		ctx := context.Background()
		todoID := int64(1)
		userID := int64(1)

		mockRepo.EXPECT().
			DeleteTodo(ctx, sqlc.DeleteTodoParams{
				ID:     todoID,
				UserID: userID,
			}).
			Return(nil)

		err := svc.DeleteTodo(ctx, todoID, userID)

		assert.NoError(t, err)
	})

	t.Run("異常系: エラーを返す", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)

		ctx := context.Background()
		todoID := int64(1)
		userID := int64(1)
		dbErr := errors.New("database error")

		mockRepo.EXPECT().
			DeleteTodo(ctx, sqlc.DeleteTodoParams{
				ID:     todoID,
				UserID: userID,
			}).
			Return(dbErr)

		err := svc.DeleteTodo(ctx, todoID, userID)

		assert.ErrorIs(t, err, dbErr)
	})
}

func TestTodoService_BatchCompleteTodos(t *testing.T) {
	t.Run("全てのTodoが存在する場合", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)
		ctx := context.Background()
		userID := int64(1)
		ids := []int64{1, 2}

		existingTodos := []sqlc.Todo{
			{ID: 1, UserID: userID, Completed: false},
			{ID: 2, UserID: userID, Completed: false},
		}
		completedTodos := []sqlc.Todo{
			{ID: 1, UserID: userID, Completed: true},
			{ID: 2, UserID: userID, Completed: true},
		}

		mockRepo.EXPECT().
			GetTodosByIDs(ctx, mock.Anything).
			Return(existingTodos, nil)
		mockRepo.EXPECT().
			BatchCompleteTodos(ctx, mock.Anything).
			Return(completedTodos, nil)

		result, err := svc.BatchCompleteTodos(ctx, userID, ids)

		assert.NoError(t, err)
		assert.Len(t, result.Succeeded, 2)
		assert.Len(t, result.Failed, 0)
	})

	t.Run("一部のIDが存在しない場合", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)
		ctx := context.Background()
		userID := int64(1)
		ids := []int64{1, 999}

		existingTodos := []sqlc.Todo{
			{ID: 1, UserID: userID},
		}
		completedTodos := []sqlc.Todo{
			{ID: 1, UserID: userID, Completed: true},
		}

		mockRepo.EXPECT().
			GetTodosByIDs(ctx, mock.Anything).
			Return(existingTodos, nil)
		mockRepo.EXPECT().
			BatchCompleteTodos(ctx, mock.Anything).
			Return(completedTodos, nil)

		result, err := svc.BatchCompleteTodos(ctx, userID, ids)

		assert.NoError(t, err)
		assert.Len(t, result.Succeeded, 1)
		assert.Len(t, result.Failed, 1)
		assert.Equal(t, int64(999), result.Failed[0].ID)
	})
}

func TestTodoService_BatchDeleteTodos(t *testing.T) {
	t.Run("全てのTodoが存在する場合", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)
		ctx := context.Background()
		userID := int64(1)
		ids := []int64{1, 2}

		existingTodos := []sqlc.Todo{
			{ID: 1, UserID: userID},
			{ID: 2, UserID: userID},
		}

		mockRepo.EXPECT().
			GetTodosByIDs(ctx, mock.Anything).
			Return(existingTodos, nil)
		mockRepo.EXPECT().
			BatchDeleteTodos(ctx, mock.Anything).
			Return(nil)

		result, err := svc.BatchDeleteTodos(ctx, userID, ids)

		assert.NoError(t, err)
		assert.Len(t, result.Succeeded, 2)
		assert.Len(t, result.Failed, 0)
		assert.ElementsMatch(t, []int64{1, 2}, result.Succeeded)
	})

	t.Run("一部のIDが存在しない場合", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)
		ctx := context.Background()
		userID := int64(1)
		ids := []int64{1, 999}

		existingTodos := []sqlc.Todo{
			{ID: 1, UserID: userID},
		}

		mockRepo.EXPECT().
			GetTodosByIDs(ctx, mock.Anything).
			Return(existingTodos, nil)
		mockRepo.EXPECT().
			BatchDeleteTodos(ctx, mock.Anything).
			Return(nil)

		result, err := svc.BatchDeleteTodos(ctx, userID, ids)

		assert.NoError(t, err)
		assert.Len(t, result.Succeeded, 1)
		assert.Equal(t, int64(1), result.Succeeded[0])
		assert.Len(t, result.Failed, 1)
		assert.Equal(t, int64(999), result.Failed[0].ID)
		assert.Equal(t, "Todo not found", result.Failed[0].Error)
	})

	t.Run("全てのIDが存在しない場合", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)
		ctx := context.Background()
		userID := int64(1)
		ids := []int64{998, 999}

		mockRepo.EXPECT().
			GetTodosByIDs(ctx, mock.Anything).
			Return([]sqlc.Todo{}, nil)

		result, err := svc.BatchDeleteTodos(ctx, userID, ids)

		assert.NoError(t, err)
		assert.Len(t, result.Succeeded, 0)
		assert.Len(t, result.Failed, 2)
	})

	t.Run("GetTodosByIDsでエラーが発生した場合", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)
		ctx := context.Background()
		userID := int64(1)
		ids := []int64{1, 2}
		dbErr := errors.New("database error")

		mockRepo.EXPECT().
			GetTodosByIDs(ctx, mock.Anything).
			Return(nil, dbErr)

		result, err := svc.BatchDeleteTodos(ctx, userID, ids)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("BatchDeleteTodosでエラーが発生した場合", func(t *testing.T) {
		mockRepo := mocks.NewMockTodoRepository(t)
		svc := NewTodoService(mockRepo)
		ctx := context.Background()
		userID := int64(1)
		ids := []int64{1, 2}
		dbErr := errors.New("database error")

		existingTodos := []sqlc.Todo{
			{ID: 1, UserID: userID},
			{ID: 2, UserID: userID},
		}

		mockRepo.EXPECT().
			GetTodosByIDs(ctx, mock.Anything).
			Return(existingTodos, nil)
		mockRepo.EXPECT().
			BatchDeleteTodos(ctx, mock.Anything).
			Return(dbErr)

		result, err := svc.BatchDeleteTodos(ctx, userID, ids)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// ヘルパー関数
func ptrString(s string) *string {
	return &s
}

func ptrBool(b bool) *bool {
	return &b
}
