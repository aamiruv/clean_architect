package usecase

import "github.com/AmirMirzayi/clean_architecture/internal/auth/domain"

type AuthService interface {
	Register(domain.Auth) error
}

type AuthUseCase struct {
	service AuthService
}

func NewAuthUseCase(authService AuthService) AuthUseCase {
	return AuthUseCase{service: authService}
}

func (u AuthUseCase) Register(auth domain.Auth) error {
	return u.service.Register(auth)
}
