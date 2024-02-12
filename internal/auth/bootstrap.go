package auth

import (
	"context"
	"database/sql"
	"github.com/AmirMirzayi/clean_architecture/api/proto/auth"
	"github.com/AmirMirzayi/clean_architecture/internal/auth/adapter/controller"
	"github.com/AmirMirzayi/clean_architecture/internal/auth/adapter/repository"
	"github.com/AmirMirzayi/clean_architecture/internal/auth/service"
	"github.com/AmirMirzayi/clean_architecture/internal/auth/usecase"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func InitializeAuthServer(server *grpc.Server, db *sql.DB) {
	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository)
	authUseCase := usecase.NewAuthUseCase(authService)
	accountHandler := controller.NewAuthHandler(authUseCase)
	auth.RegisterAuthServiceServer(server, accountHandler)
}

func RegisterGateway(ctx context.Context, mux *runtime.ServeMux, grpcAddress string, options ...grpc.DialOption) error {
	if err := auth.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAddress, options); err != nil {
		return err
	}
	return nil
}
