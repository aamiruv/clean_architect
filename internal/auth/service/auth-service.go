package service

import "github.com/AmirMirzayi/clean_architecture/internal/auth/domain"

type AuthRepository interface {
	Register(domain.Auth) error
}

type AuthService struct {
	repository AuthRepository
}

func NewAuthService(authRepository AuthRepository) AuthService {
	return AuthService{repository: authRepository}
}

func (s AuthService) Register(auth domain.Auth) error {
	return s.repository.Register(auth)
}
