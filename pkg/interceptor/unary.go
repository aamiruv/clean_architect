// Package interceptors contains grpc interceptors which acts before request get into the server handler
package interceptor

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Recovery(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, logger *log.Logger) (_ any, err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("panic recovered: %v, in: %s", r, info.FullMethod)
			err = status.Error(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}

func ResponseTimeMeter(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, logger *log.Logger) (resp any, err error) {
	now := time.Now()
	defer func() {
		logger.Printf("completed in: %v", time.Since(now))
	}()
	return handler(ctx, req)
}

// DenyUnauthorizedClient will check request for authorization metadata
// and returns error if not found
func DenyUnauthorizedClient(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// fixme: validate authorization token
	v := metadata.ValueFromIncomingContext(ctx, "authorization")
	if len(v) == 0 {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}
	return handler(ctx, req)
}
