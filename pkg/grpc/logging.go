package grpc

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func UnaryServerLoggerInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = logger.WithContext(ctx)
		return handler(ctx, req)
	}
}
