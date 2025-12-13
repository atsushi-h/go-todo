package service

import (
	"context"
	"errors"
	"fmt"

	"go-todo/db/sqlc"
	"go-todo/internal/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/markbates/goth"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	repo      UserRepository
	pool      *pgxpool.Pool
	txManager database.TxManager
}

func NewUserService(repo UserRepository, pool *pgxpool.Pool) *UserService {
	return &UserService{
		repo:      repo,
		pool:      pool,
		txManager: database.NewTxManager(pool),
	}
}

func (s *UserService) FindOrCreateFromOAuth(ctx context.Context, gothUser goth.User) (*sqlc.User, error) {
	// 既存ユーザーを検索
	user, err := s.repo.GetUserByProviderID(ctx, sqlc.GetUserByProviderIDParams{
		Provider:   gothUser.Provider,
		ProviderID: gothUser.UserID,
	})

	if err == nil {
		// 既存ユーザーの情報を更新
		var avatarURL *string
		if gothUser.AvatarURL != "" {
			avatarURL = &gothUser.AvatarURL
		}
		updated, err := s.repo.UpdateUser(ctx, sqlc.UpdateUserParams{
			ID:        user.ID,
			Name:      gothUser.Name,
			AvatarUrl: avatarURL,
		})
		if err != nil {
			return nil, err
		}
		return &updated, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// 新規ユーザー作成
	var avatarURL *string
	if gothUser.AvatarURL != "" {
		avatarURL = &gothUser.AvatarURL
	}
	newUser, err := s.repo.CreateUser(ctx, sqlc.CreateUserParams{
		Email:      gothUser.Email,
		Name:       gothUser.Name,
		AvatarUrl:  avatarURL,
		Provider:   gothUser.Provider,
		ProviderID: gothUser.UserID,
	})
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*sqlc.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) DeleteAccount(ctx context.Context, userID int64) error {
	// ユーザー存在確認
	_, err := s.repo.GetUserByID(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	// トランザクション内で削除実行
	return s.txManager.RunInTx(ctx, func(tx pgx.Tx) error {
		queries := sqlc.New(tx)

		// Step 1: ユーザーの全Todoをソフトデリート
		if err := queries.DeleteTodosByUserID(ctx, userID); err != nil {
			return fmt.Errorf("delete todos: %w", err)
		}

		// Step 2: ユーザーをソフトデリート
		if err := queries.DeleteUser(ctx, userID); err != nil {
			return fmt.Errorf("delete user: %w", err)
		}

		return nil
	})
}
