package service

import (
	"context"
	"errors"

	"github.com/amirzayi/clean_architect/internal/domain"
	"github.com/amirzayi/clean_architect/pkg/auth"
	"github.com/amirzayi/clean_architect/pkg/hash"
)

type Auth interface {
	Register(ctx context.Context, auth domain.Auth) error
	Login(ctx context.Context, auth domain.Auth) (token string, err error)
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
		Role:        domain.UserRoleNormal,
	}

	_, err = a.userService.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
func (a *authService) Login(ctx context.Context, auth domain.Auth) (string, error) {
	// todo: add login via phone no
	user, err := a.userService.GetByEmail(ctx, auth.Email)
	if err != nil {
		return "", err
	}

	if user.Status == domain.UserStatusBanned {
		return "", errors.New("user is banned")
	}

	if err = a.hasher.Compare(user.Password, auth.Password); err != nil {
		return "", err
	}

	token, err := a.authManager.CreateToken(user.ID, string(user.Role))
	return token, err
}
