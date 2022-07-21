package jwt

import (
	"context"
	"errors"
	"strings"

	grpc "github.com/ppcamp/go-grpc"
)

const (
	jwtAuthFieldName          = "authorization"
	ContextKey       grpc.Key = "jwt"
)

var (
	ErrMissingBearer error = errors.New("missing token")
	ErrBadAuthString error = errors.New("bad authorization string")
)

type AuthenticationService interface {
	// Data return the parsed jwt token data field, usually known as session.
	Data(signedToken string) (any, error)

	// Allow is used to check if the current path is an unprotected endpoint. If so, allow it.
	Allow(path string) bool
}

// jwtFromHeader tries to fetch the jwt bearer token from the metadata gRPC field.
func jwtFromHeader(ctx context.Context) (string, error) {
	v, err := grpc.MetadataField(ctx, jwtAuthFieldName)
	if err != nil {
		return "", err
	}

	val, ok := v.(string)
	if !ok {
		return "", errors.New("fail to cast to string")
	}

	if val == "" {
		return "", ErrMissingBearer
	}

	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", ErrBadAuthString
	}

	return splits[1], nil
}
