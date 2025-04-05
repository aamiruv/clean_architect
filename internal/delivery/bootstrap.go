package delivery

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	grpcapi "github.com/amirzayi/clean_architect/api/grpc"
	"github.com/amirzayi/clean_architect/api/http/handler"
	"github.com/amirzayi/clean_architect/api/proto/authpb"
	"github.com/amirzayi/clean_architect/internal/service"
)

func SetupHTTPRouter(mux *http.ServeMux, logger *log.Logger, services *service.Services) {
	handler.Register(mux, logger, services)
}

func SetupGRPC(server *grpc.Server, services *service.Services) {
	authService := grpcapi.NewAuthGrpcService(services.Auth)
	authpb.RegisterAuthServiceServer(server, authService)
}

func SetupGRPCGateway(ctx context.Context, grpcAddress string, mux *runtime.ServeMux, options ...grpc.DialOption) error {
	if err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAddress, options); err != nil {
		return err
	}

	return nil
}
