package log

import (
	"context"
	"time"

	"github.com/google/uuid"
	w "github.com/ppcamp/go-grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

const LogId w.Key = "id"

func logTime(logger grpclog.LoggerV2, id string, s time.Time) {
	elapsedMs := time.Since(s).Milliseconds()
	logger.Infof("request %s took %dms", id, elapsedMs)
}

// UnaryInterceptor will use the log param to log the request duration. It also adds the request's
// id to the context
func UnaryInterceptor(log grpclog.LoggerV2) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler,
	) (resp interface{}, err error) {
		s := time.Now()

		// generating a unique id
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "fail to create request id: %w", err)
		}
		defer logTime(log, id.String(), s)

		// adding the id to the context
		ctx = context.WithValue(ctx, LogId, id)
		return next(ctx, req)
	}
}

// StreamInterceptor will use the log param to log the request duration. It also adds the request's
// id to the context
func StreamInterceptor(log grpclog.LoggerV2) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		next grpc.StreamHandler,
	) error {
		s := time.Now()

		// generating a unique id
		id, err := uuid.NewRandom()
		if err != nil {
			return status.Errorf(codes.Internal, "fail to create request id: %w", err)
		}
		defer logTime(log, id.String(), s)

		// adding the id to the context
		ctx := context.WithValue(ss.Context(), LogId, id)
		ss = &w.ServerStream{ServerStream: ss, Ctx: ctx}

		return next(srv, ss)
	}
}
