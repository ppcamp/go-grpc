package ratelimit

import (
	"context"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverClose interface{ GracefulStop() }

// UnaryInterceptor implements a gRPC middleware that allows n requests before starts rejecting.
// The interceptor is threadsafe (uses atomic operations).
func UnaryInterceptor(n uint64) grpc.UnaryServerInterceptor {
	c := uint64(0)

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler,
	) (resp interface{}, err error) {
		atomic.AddUint64(&c, 1)
		if atomic.LoadUint64(&c) == n {
			s, ok := info.Server.(serverClose)
			if !ok {
				panic("fail to get the server")
			}
			defer s.GracefulStop()

			return nil, status.Errorf(codes.ResourceExhausted, "%s was rejected by ratelimit", info.FullMethod)
		}
		return next(ctx, req)
	}
}

// StreamInterceptor implements a gRPC middleware that allows n requests before starts rejecting.
// The interceptor is threadsafe (uses atomic operations).
func StreamInterceptor(n uint64) grpc.StreamServerInterceptor {
	c := uint64(0)

	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		next grpc.StreamHandler,
	) error {
		atomic.AddUint64(&c, 1)
		if atomic.LoadUint64(&c) == n {
			s, ok := srv.(serverClose)
			if !ok {
				panic("fail to get the server")
			}
			defer s.GracefulStop()

			return status.Errorf(codes.ResourceExhausted, "%s was rejected by ratelimit", info.FullMethod)
		}
		return next(srv, ss)
	}
}
