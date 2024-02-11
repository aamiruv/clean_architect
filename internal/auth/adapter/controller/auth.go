package controller

import (
	"context"
	"github.com/AmirMirzayi/clean_architecture/api/proto/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterGateway(ctx context.Context, mux *runtime.ServeMux, grpcAddress string, options ...grpc.DialOption) error {
	if err := auth.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAddress, options); err != nil {
		return err
	}
	return nil
}

type AuthUseCase interface {
}

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	useCase AuthUseCase
}

func NewAuthHandler(authUseCase AuthUseCase) AuthHandler {
	return AuthHandler{useCase: authUseCase}
}

func (h AuthHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	return &auth.RegisterResponse{UserId: "user" + req.GetUserName() + req.GetPassword()}, nil
}
