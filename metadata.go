package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"
)

var (
	ErrFailRetrieveMetadata error = errors.New("failed to get metadata from context")
	ErrFieldNotFound        error = errors.New("field wasn't found")
)

func MetadataField(ctx context.Context, field string) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrFailRetrieveMetadata
	}

	kkeys := md.Get(field)
	if len(kkeys) == 0 {
		return nil, ErrFieldNotFound
	}

	return kkeys[0], nil
}
