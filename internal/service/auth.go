package service

import (
	"context"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/pkg/auth"
	"github.com/amirzayi/clean_architect/pkg/hash"
)

type Auth interface {
	Register(ctx context.Context, auth domain.Auth) error
}

type authService struct {
	userService User
	hasher      hash.PasswordHasher
	authManager auth.Manager
}

func NewAuthService(userService User, hasher hash.PasswordHasher, authManager auth.Manager) Auth {
	return &authService{
		userService: userService,
		hasher:      hasher,
		authManager: authManager,
	}
}

func (a *authService) Register(ctx context.Context, auth domain.Auth) error {
	pwd, err := a.hasher.Hash(auth.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Email:       auth.Email,
		PhoneNumber: auth.PhoneNumber,
		Password:    pwd,
	}

	_, err = a.userService.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
