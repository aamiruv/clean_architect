// Package usecase is responsible for some operation to be executed.
package usecase

import (
	"context"

	"github.com/amirzayi/clean_architec/internal/auth/domain"
	userDomain "github.com/amirzayi/clean_architec/internal/user/domain"
	"github.com/amirzayi/clean_architec/internal/user/service"
)

type AuthService interface {
	HashPassword(context.Context, string) (string, error)
}

type AuthUseCase struct {
	authService AuthService
	userService service.UserService
}

func NewAuthUseCase(authService AuthService, userService service.UserService) AuthUseCase {
	return AuthUseCase{
		authService: authService,
		userService: userService,
	}
}

func (u AuthUseCase) Register(ctx context.Context, auth domain.Auth) error {
	pwd, err := u.authService.HashPassword(ctx, auth.Password)
	if err != nil {
		return err
	}

	user := userDomain.User{
		Email:       auth.Email,
		PhoneNumber: auth.PhoneNumber,
		Password:    pwd,
	}
	_, err = u.userService.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
