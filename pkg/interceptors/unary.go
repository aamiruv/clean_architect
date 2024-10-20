package interceptors

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Recovery(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered: %v, in: %s", r, info.FullMethod)
			err = status.Error(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}

func ResponseTimeMeter(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	now := time.Now()
	defer func() {
		log.Printf("completed in: %v", time.Since(now))
	}()
	return handler(ctx, req)
}
