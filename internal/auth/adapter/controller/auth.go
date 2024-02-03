package controller

import (
	"context"
	"github.com/AmirMirzayi/clean_architecture/api/proto/auth"
)

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
}

func NewAuthHandler() AuthHandler {
	return AuthHandler{}
}

func (h AuthHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	return &auth.RegisterResponse{UserId: "user" + req.GetUserName() + req.GetPassword()}, nil
}
