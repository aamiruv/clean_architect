package grpc

import (
		"context"

		"google.golang.org/grpc/codes"
		"google.golang.org/grpc/status"
		"google.golang.org/protobuf/types/known/emptypb"

		"github.com/amirzayi/clean_architect/api/proto/authpb"
		"github.com/amirzayi/clean_architect/internal/domain"
		"github.com/amirzayi/clean_architect/internal/service"
)

type authService struct {
		authpb.UnimplementedAuthServiceServer
		auth service.Auth
}

func NewAuthGrpcService(auth service.Auth) *authService {
		return &authService{auth: auth}
}

func (h *authService) Register(ctx context.Context, req *authpb.RegisterRequest) (*emptypb.Empty, error) {
		auth := domain.Auth{
				Email:       req.GetEmail(),
				PhoneNumber: req.GetPhoneNumber(),
				Password:    req.GetPassword(),
		}
		if err := h.auth.Register(ctx, auth); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to register: %v", err)
		}
		return &emptypb.Empty{}, nil
}
