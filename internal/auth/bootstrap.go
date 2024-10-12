// Package auth holds all bootstrapper that manages business flow related to user authentication & authorization.
package auth

import (
	"context"
	"database/sql"

	"github.com/AmirMirzayi/clean_architecture/api/proto/authpb"
	"github.com/AmirMirzayi/clean_architecture/internal/auth/adapter/controller"
	"github.com/AmirMirzayi/clean_architecture/internal/auth/service"
	"github.com/AmirMirzayi/clean_architecture/internal/auth/usecase"
	"github.com/AmirMirzayi/clean_architecture/internal/user"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func InitializeAuthServer(server *grpc.Server, db *sql.DB) {
	userService := user.NewService(user.NewSQLRepository(db))

	authService := service.NewAuthService()
	authUseCase := usecase.NewAuthUseCase(authService, userService)
	accountHandler := controller.NewAuthHandler(authUseCase)
	authpb.RegisterAuthServiceServer(server, accountHandler)
}

func RegisterGateway(ctx context.Context, mux *runtime.ServeMux, grpcAddress string, options ...grpc.DialOption) error {
	if err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAddress, options); err != nil {
		return err
	}
	return nil
}
