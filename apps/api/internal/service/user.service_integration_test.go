//go:build integration

package service

import (
	"context"
	"testing"
	"time"

	"go-todo/db/sqlc"
	"go-todo/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_DeleteAccount_Integration(t *testing.T) {
	ctx := context.Background()

	pool, err := database.NewPool(ctx)
	require.NoError(t, err)
	defer pool.Close()

	queries := sqlc.New(pool)
	userService := NewUserService(queries, pool)

	t.Run("正常系: ユーザーとTodoが削除される", func(t *testing.T) {
		// テストユーザー作成
		user, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
			Email:      "delete-test@example.com",
			Name:       "Delete Test",
			Provider:   "test",
			ProviderID: "test-delete-" + time.Now().Format("20060102150405"),
		})
		require.NoError(t, err)

		// テストTodo作成
		_, err = queries.CreateTodo(ctx, sqlc.CreateTodoParams{
			UserID:      user.ID,
			Title:       "Todo 1",
			Description: ptrString("Description 1"),
		})
		require.NoError(t, err)

		// 削除実行
		err = userService.DeleteAccount(ctx, user.ID)
		require.NoError(t, err)

		// 検証: ユーザーが取得できない
		_, err = queries.GetUserByID(ctx, user.ID)
		assert.Error(t, err)

		// 検証: Todoが取得できない
		todos, err := queries.ListTodosByUser(ctx, user.ID)
		require.NoError(t, err)
		assert.Empty(t, todos)

		// 検証: ソフトデリート確認（物理削除されていない）
		var count int
		err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE id = $1", user.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "User should still exist in database")
	})
}
