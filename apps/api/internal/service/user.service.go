package service

import (
	"context"

	"go-todo/internal/model"
	"go-todo/internal/repository"

	"github.com/markbates/goth"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// OAuthユーザーを検索または作成
func (s *UserService) FindOrCreateFromOAuth(ctx context.Context, gothUser goth.User) (*model.User, error) {
	// 既存ユーザーを検索
	user, err := s.repo.FindByProviderID(ctx, gothUser.Provider, gothUser.UserID)
	if err == nil {
		// 既存ユーザーの情報を更新
		user.Name = gothUser.Name
		user.AvatarURL = gothUser.AvatarURL
		if err := s.repo.Update(ctx, user); err != nil {
			return nil, err
		}
		return user, nil
	}

	if err != repository.ErrUserNotFound {
		return nil, err
	}

	// 新規ユーザー作成
	user = &model.User{
		Email:      gothUser.Email,
		Name:       gothUser.Name,
		AvatarURL:  gothUser.AvatarURL,
		Provider:   gothUser.Provider,
		ProviderID: gothUser.UserID,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}
