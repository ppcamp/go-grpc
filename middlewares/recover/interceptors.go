package recover

import (
	"context"
	"errors"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// UnaryInterceptor implements a gRPC middleware that allows the server to recover from a gRPC panic
func UnaryInterceptor(log grpclog.LoggerV2) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			panicked := recover()
			if panicked != nil {
				err = errors.New("panicked, check out server logs")
				log.Errorf("recovered from a panic, reason: %v", panicked)
				debug.PrintStack()
			}
			if err != nil {
				log.Errorf("some error ocurred, err: %v", err)
			}
		}()
		return next(ctx, req)
	}
}

// StreamInterceptor implements a gRPC middleware that allows the server to recover from a gRPC panic
func StreamInterceptor(log grpclog.LoggerV2) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		next grpc.StreamHandler,
	) (err error) {
		defer func() {
			panicked := recover()
			if panicked != nil {
				err = errors.New("panicked, check out server logs")
				log.Errorf("recovered from a panic, reason: %v", panicked)
				debug.PrintStack()
			}
			if err != nil {
				log.Errorf("some error ocurred, err: %v", err)
			}
		}()
		return next(srv, ss)
	}
}
