package grpc

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// SecretValidator provides gRPC interceptors for validating client secret.
type SecretValidator struct {
	SecretToken string
}

// UnaryInterceptor is an gRPC Interceptor implementing client token validation.
func (s *SecretValidator) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
	}
	if len(meta["token"]) != 1 {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}
	if meta["token"][0] != s.SecretToken {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}
