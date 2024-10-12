package service

import (
	"context"
)

type AuthService struct{}

func NewAuthService() AuthService {
	return AuthService{}
}

func (s AuthService) HashPassword(ctx context.Context, pwd string) (string, error) {
	// todo: provide some password encryption mechanism
	return pwd, nil
}
