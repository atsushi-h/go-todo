package service

import (
	"context"
	"errors"

	"go-todo/db/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/markbates/goth"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
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
