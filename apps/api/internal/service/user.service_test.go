package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-todo/db/sqlc"
	"go-todo/internal/service/mocks"

	"github.com/jackc/pgx/v5"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetByID(t *testing.T) {
	t.Run("正常系: Userを取得できる", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		svc := NewUserService(mockRepo, nil)

		ctx := context.Background()
		userID := int64(1)
		now := time.Now()

		expectedUser := sqlc.User{
			ID:         userID,
			Email:      "test@example.com",
			Name:       "Test User",
			AvatarUrl:  ptrString("https://example.com/avatar.png"),
			Provider:   "google",
			ProviderID: "google-123",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		mockRepo.EXPECT().
			GetUserByID(ctx, userID).
			Return(expectedUser, nil)

		result, err := svc.GetByID(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, expectedUser.ID, result.ID)
		assert.Equal(t, expectedUser.Email, result.Email)
		assert.Equal(t, expectedUser.Name, result.Name)
	})

	t.Run("異常系: ErrUserNotFoundを返す", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		svc := NewUserService(mockRepo, nil)

		ctx := context.Background()
		userID := int64(999)

		mockRepo.EXPECT().
			GetUserByID(ctx, userID).
			Return(sqlc.User{}, pgx.ErrNoRows)

		result, err := svc.GetByID(ctx, userID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("異常系: その他のエラーをそのまま返す", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		svc := NewUserService(mockRepo, nil)

		ctx := context.Background()
		userID := int64(1)
		dbErr := errors.New("database error")

		mockRepo.EXPECT().
			GetUserByID(ctx, userID).
			Return(sqlc.User{}, dbErr)

		result, err := svc.GetByID(ctx, userID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestUserService_FindOrCreateFromOAuth(t *testing.T) {
	t.Run("正常系: 既存ユーザーを更新して返す", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		svc := NewUserService(mockRepo, nil)

		ctx := context.Background()
		now := time.Now()

		gothUser := goth.User{
			Provider:  "google",
			UserID:    "google-123",
			Email:     "test@example.com",
			Name:      "Updated Name",
			AvatarURL: "https://example.com/new-avatar.png",
		}

		existingUser := sqlc.User{
			ID:         1,
			Email:      "test@example.com",
			Name:       "Old Name",
			AvatarUrl:  ptrString("https://example.com/old-avatar.png"),
			Provider:   "google",
			ProviderID: "google-123",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		updatedUser := sqlc.User{
			ID:         1,
			Email:      "test@example.com",
			Name:       "Updated Name",
			AvatarUrl:  ptrString("https://example.com/new-avatar.png"),
			Provider:   "google",
			ProviderID: "google-123",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		mockRepo.EXPECT().
			GetUserByProviderID(ctx, sqlc.GetUserByProviderIDParams{
				Provider:   "google",
				ProviderID: "google-123",
			}).
			Return(existingUser, nil)

		mockRepo.EXPECT().
			UpdateUser(ctx, sqlc.UpdateUserParams{
				ID:        existingUser.ID,
				Name:      gothUser.Name,
				AvatarUrl: ptrString(gothUser.AvatarURL),
			}).
			Return(updatedUser, nil)

		result, err := svc.FindOrCreateFromOAuth(ctx, gothUser)

		require.NoError(t, err)
		assert.Equal(t, updatedUser.ID, result.ID)
		assert.Equal(t, updatedUser.Name, result.Name)
		assert.Equal(t, updatedUser.AvatarUrl, result.AvatarUrl)
	})

	t.Run("正常系: 新規ユーザーを作成して返す", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		svc := NewUserService(mockRepo, nil)

		ctx := context.Background()
		now := time.Now()

		gothUser := goth.User{
			Provider:  "google",
			UserID:    "google-456",
			Email:     "new@example.com",
			Name:      "New User",
			AvatarURL: "https://example.com/avatar.png",
		}

		newUser := sqlc.User{
			ID:         2,
			Email:      "new@example.com",
			Name:       "New User",
			AvatarUrl:  ptrString("https://example.com/avatar.png"),
			Provider:   "google",
			ProviderID: "google-456",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		mockRepo.EXPECT().
			GetUserByProviderID(ctx, sqlc.GetUserByProviderIDParams{
				Provider:   "google",
				ProviderID: "google-456",
			}).
			Return(sqlc.User{}, pgx.ErrNoRows)

		mockRepo.EXPECT().
			CreateUser(ctx, sqlc.CreateUserParams{
				Email:      gothUser.Email,
				Name:       gothUser.Name,
				AvatarUrl:  ptrString(gothUser.AvatarURL),
				Provider:   gothUser.Provider,
				ProviderID: gothUser.UserID,
			}).
			Return(newUser, nil)

		result, err := svc.FindOrCreateFromOAuth(ctx, gothUser)

		require.NoError(t, err)
		assert.Equal(t, newUser.ID, result.ID)
		assert.Equal(t, newUser.Email, result.Email)
		assert.Equal(t, newUser.Name, result.Name)
	})

	t.Run("正常系: AvatarURLが空の場合はnilで作成", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		svc := NewUserService(mockRepo, nil)

		ctx := context.Background()
		now := time.Now()

		gothUser := goth.User{
			Provider:  "google",
			UserID:    "google-789",
			Email:     "noavatar@example.com",
			Name:      "No Avatar User",
			AvatarURL: "",
		}

		newUser := sqlc.User{
			ID:         3,
			Email:      "noavatar@example.com",
			Name:       "No Avatar User",
			AvatarUrl:  nil,
			Provider:   "google",
			ProviderID: "google-789",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		mockRepo.EXPECT().
			GetUserByProviderID(ctx, sqlc.GetUserByProviderIDParams{
				Provider:   "google",
				ProviderID: "google-789",
			}).
			Return(sqlc.User{}, pgx.ErrNoRows)

		mockRepo.EXPECT().
			CreateUser(ctx, sqlc.CreateUserParams{
				Email:      gothUser.Email,
				Name:       gothUser.Name,
				AvatarUrl:  nil,
				Provider:   gothUser.Provider,
				ProviderID: gothUser.UserID,
			}).
			Return(newUser, nil)

		result, err := svc.FindOrCreateFromOAuth(ctx, gothUser)

		require.NoError(t, err)
		assert.Equal(t, newUser.ID, result.ID)
		assert.Nil(t, result.AvatarUrl)
	})

	t.Run("異常系: GetUserByProviderIDでエラー", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		svc := NewUserService(mockRepo, nil)

		ctx := context.Background()
		dbErr := errors.New("database error")

		gothUser := goth.User{
			Provider: "google",
			UserID:   "google-123",
		}

		mockRepo.EXPECT().
			GetUserByProviderID(ctx, sqlc.GetUserByProviderIDParams{
				Provider:   "google",
				ProviderID: "google-123",
			}).
			Return(sqlc.User{}, dbErr)

		result, err := svc.FindOrCreateFromOAuth(ctx, gothUser)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("異常系: UpdateUserでエラー", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		svc := NewUserService(mockRepo, nil)

		ctx := context.Background()
		now := time.Now()
		dbErr := errors.New("update error")

		gothUser := goth.User{
			Provider:  "google",
			UserID:    "google-123",
			Name:      "Updated Name",
			AvatarURL: "https://example.com/avatar.png",
		}

		existingUser := sqlc.User{
			ID:         1,
			Email:      "test@example.com",
			Name:       "Old Name",
			Provider:   "google",
			ProviderID: "google-123",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		mockRepo.EXPECT().
			GetUserByProviderID(ctx, sqlc.GetUserByProviderIDParams{
				Provider:   "google",
				ProviderID: "google-123",
			}).
			Return(existingUser, nil)

		mockRepo.EXPECT().
			UpdateUser(ctx, sqlc.UpdateUserParams{
				ID:        existingUser.ID,
				Name:      gothUser.Name,
				AvatarUrl: ptrString(gothUser.AvatarURL),
			}).
			Return(sqlc.User{}, dbErr)

		result, err := svc.FindOrCreateFromOAuth(ctx, gothUser)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("異常系: CreateUserでエラー", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		svc := NewUserService(mockRepo, nil)

		ctx := context.Background()
		dbErr := errors.New("create error")

		gothUser := goth.User{
			Provider:  "google",
			UserID:    "google-new",
			Email:     "new@example.com",
			Name:      "New User",
			AvatarURL: "https://example.com/avatar.png",
		}

		mockRepo.EXPECT().
			GetUserByProviderID(ctx, sqlc.GetUserByProviderIDParams{
				Provider:   "google",
				ProviderID: "google-new",
			}).
			Return(sqlc.User{}, pgx.ErrNoRows)

		mockRepo.EXPECT().
			CreateUser(ctx, sqlc.CreateUserParams{
				Email:      gothUser.Email,
				Name:       gothUser.Name,
				AvatarUrl:  ptrString(gothUser.AvatarURL),
				Provider:   gothUser.Provider,
				ProviderID: gothUser.UserID,
			}).
			Return(sqlc.User{}, dbErr)

		result, err := svc.FindOrCreateFromOAuth(ctx, gothUser)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestUserService_DeleteAccount(t *testing.T) {
	t.Run("異常系: ユーザーが存在しない場合", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		svc := NewUserService(mockRepo, nil)

		ctx := context.Background()
		userID := int64(999)

		mockRepo.EXPECT().
			GetUserByID(ctx, userID).
			Return(sqlc.User{}, pgx.ErrNoRows)

		err := svc.DeleteAccount(ctx, userID)

		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("異常系: GetUserByIDでその他のエラー", func(t *testing.T) {
		mockRepo := mocks.NewMockUserRepository(t)
		svc := NewUserService(mockRepo, nil)

		ctx := context.Background()
		userID := int64(1)
		dbErr := errors.New("database error")

		mockRepo.EXPECT().
			GetUserByID(ctx, userID).
			Return(sqlc.User{}, dbErr)

		err := svc.DeleteAccount(ctx, userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get user")
	})
}
