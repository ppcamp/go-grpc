package grpc

import (
	"context"

	"google.golang.org/grpc"
)

type ServerStream struct {
	grpc.ServerStream

	Ctx context.Context
}

func (s *ServerStream) Context() context.Context {
	return s.Ctx
}
