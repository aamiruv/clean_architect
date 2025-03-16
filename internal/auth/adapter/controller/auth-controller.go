// Package controller holds grpc service implementation.
package controller

import (
	"context"

	"github.com/amirzayi/clean_architec/api/proto/authpb"
	"github.com/amirzayi/clean_architec/internal/auth/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthUseCase interface {
	Register(context.Context, domain.Auth) error
}

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	useCase AuthUseCase
}

func NewAuthGRPCHandler(authUseCase AuthUseCase) *AuthHandler {
	return &AuthHandler{useCase: authUseCase}
}

func (h AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*emptypb.Empty, error) {
	auth := domain.Auth{
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
		Password:    req.GetPassword(),
	}
	if err := h.useCase.Register(ctx, auth); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register: %v", err)
	}
	return &emptypb.Empty{}, nil
}
