package controller

import (
	"context"
	"errors"
	"github.com/AmirMirzayi/clean_architecture/api/proto/auth"
	"github.com/AmirMirzayi/clean_architecture/internal/auth/domain"
)

type AuthUseCase interface {
	Register(domain.Auth) error
}

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	useCase AuthUseCase
}

func NewAuthHandler(authUseCase AuthUseCase) AuthHandler {
	return AuthHandler{useCase: authUseCase}
}

func (h AuthHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	d := domain.Auth{UserName: req.GetUserName(), Password: req.GetPassword()}
	if h.useCase.Register(d) != nil {
		return nil, errors.New("something went wrong")
	}
	return &auth.RegisterResponse{UserId: "user" + req.GetUserName() + req.GetPassword()}, nil
}
