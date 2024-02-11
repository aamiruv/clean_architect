package usecase

type AuthService interface {
}

type AuthUseCase struct {
	service AuthService
}

func NewAuthUseCase(authService AuthService) AuthUseCase {
	return AuthUseCase{service: authService}
}
