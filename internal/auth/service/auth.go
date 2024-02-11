package service

type AuthRepository interface {
}

type AuthService struct {
	repository AuthRepository
}

func NewAuthService(authRepository AuthRepository) AuthService {
	return AuthService{repository: authRepository}
}
