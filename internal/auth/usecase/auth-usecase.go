package usecase

import (
	"context"

	"github.com/AmirMirzayi/clean_architecture/internal/auth/domain"
	userDomain "github.com/AmirMirzayi/clean_architecture/internal/user/domain"
	"github.com/AmirMirzayi/clean_architecture/internal/user/service"
)

type AuthService interface {
	Register(context.Context, domain.Auth) (domain.Auth, error)
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
	auth, err := u.authService.Register(ctx, auth)
	if err != nil {
		return err
	}

	user := userDomain.User{
		Email:       auth.Email,
		PhoneNumber: auth.PhoneNumber,
		Password:    auth.Password,
	}
	_, err = u.userService.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
