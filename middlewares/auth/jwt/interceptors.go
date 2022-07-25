package jwt

import (
	"context"
	"errors"

	wrappers "github.com/ppcamp/go-grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryInterceptor will use the authentication service to get the data from JWT token, and put it
// into a new context, that will be passed to the interceptors chain
func UnaryInterceptor(aservice AuthenticationService) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler,
	) (resp interface{}, err error) {

		token, err := jwtFromHeader(ctx)

		if errors.Is(err, wrappers.ErrFieldNotFound) {
			// check if the endpoint is protected or not
			if !aservice.Allow(info.FullMethod) {
				return nil, status.Errorf(codes.Unauthenticated, "protected endpoint: %w", err)
			}
			ctx = context.WithValue(ctx, ContextKey, nil)
			return next(ctx, req)
		}

		if err != nil {
			return nil, status.Errorf(codes.Internal, "fail when getting jwt token: %w", err)
		}

		parsed, err := aservice.Data(token)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "fail parse data from token: %w", err)
		}

		ctx = context.WithValue(ctx, ContextKey, parsed)
		return next(ctx, req)
	}
}

// StreamInterceptor will use the authentication service to get the data from JWT token, and put it
// into a new context, that will be passed to the interceptors chain
func StreamInterceptor(aservice AuthenticationService) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		next grpc.StreamHandler,
	) error {
		ctx := ss.Context()

		token, err := jwtFromHeader(ctx)

		if errors.Is(err, wrappers.ErrFieldNotFound) {
			// check if the endpoint is protected or not
			if !aservice.Allow(info.FullMethod) {
				return status.Errorf(codes.Unauthenticated, "protected endpoint: %w", err)
			}
			ctx = context.WithValue(ctx, ContextKey, nil)
			ss = &wrappers.ServerStream{ServerStream: ss, Ctx: ctx}
			return next(srv, ss)
		}

		if err != nil {
			return status.Errorf(codes.Internal, "fail when getting jwt token: %w", err)
		}

		parsed, err := aservice.Data(token)
		if err != nil {
			return status.Errorf(codes.Internal, "fail parse data from token: %w", err)
		}

		ctx = context.WithValue(ctx, ContextKey, parsed)
		ss = &wrappers.ServerStream{ServerStream: ss, Ctx: ctx}
		return next(srv, ss)
	}
}
