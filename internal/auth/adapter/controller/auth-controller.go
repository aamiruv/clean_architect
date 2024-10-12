package controller

import (
	"context"

	"github.com/AmirMirzayi/clean_architecture/api/proto/authpb"
	"github.com/AmirMirzayi/clean_architecture/internal/auth/domain"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthUseCase interface {
	Register(context.Context, domain.Auth) error
}

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	useCase AuthUseCase
}

func NewAuthHandler(authUseCase AuthUseCase) AuthHandler {
	return AuthHandler{useCase: authUseCase}
}

func (h AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*emptypb.Empty, error) {
	auth := domain.Auth{
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
		Password:    req.GetPassword(),
	}
	if err := h.useCase.Register(ctx, auth); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
