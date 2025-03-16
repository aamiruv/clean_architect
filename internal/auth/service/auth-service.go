package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func NewAuthService() AuthService {
	return AuthService{}
}

func (s AuthService) HashPassword(ctx context.Context, pwd string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(password), nil
}
