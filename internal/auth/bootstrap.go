// Package auth holds all bootstrapper that manages business flow related to user authentication & authorization.
package auth

import (
	"context"
	"database/sql"

	"github.com/amirzayi/clean_architec/api/proto/authpb"
	"github.com/amirzayi/clean_architec/internal/auth/adapter/controller"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func InitializeAuthServer(server *grpc.Server, db *sql.DB, authUseCase controller.AuthUseCase) {
	accountHandler := controller.NewAuthGRPCHandler(authUseCase)
	authpb.RegisterAuthServiceServer(server, accountHandler)
}

func RegisterGateway(ctx context.Context, mux *runtime.ServeMux, grpcAddress string, options ...grpc.DialOption) error {
	if err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAddress, options); err != nil {
		return err
	}
	return nil
}
