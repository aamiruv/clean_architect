package service

import (
	"context"

	"github.com/AmirMirzayi/clean_architecture/internal/auth/domain"
)

type AuthService struct{}

func NewAuthService() AuthService {
	return AuthService{}
}

func (s AuthService) Register(ctx context.Context, auth domain.Auth) (domain.Auth, error) {
	// todo: provide some password encryption mechanism
	return domain.Auth{
		Email:       auth.Email,
		PhoneNumber: auth.PhoneNumber,
		Password:    auth.Password,
	}, nil
}
