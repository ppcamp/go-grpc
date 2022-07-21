package ratelimit

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverClose interface {
	GracefulStop()
}

func UnaryInterceptor(n int) grpc.UnaryServerInterceptor {
	c := 0
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler,
	) (resp interface{}, err error) {
		c++
		if c == n {
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

func StreamInterceptor(n int) grpc.StreamServerInterceptor {
	c := 0
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		next grpc.StreamHandler,
	) error {
		c++
		if c == n {
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
