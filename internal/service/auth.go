package service

import (
	"context"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/pkg/hash"
)

type Auth interface {
	Register(ctx context.Context, auth domain.Auth) error
}

type auth struct {
	userService User
	hasher      hash.PasswordHasher
}

func NewAuthService(userService User, hasher hash.PasswordHasher) Auth {
	return &auth{
		userService: userService,
		hasher:      hasher,
	}
}

func (a *auth) Register(ctx context.Context, auth domain.Auth) error {
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
