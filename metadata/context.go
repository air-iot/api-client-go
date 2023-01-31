package metadata

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func GetGrpcContext(ctx context.Context, data map[string]string) context.Context {
	md := metadata.New(data)
	// 发送 metadata
	// 创建带有meta的context
	return metadata.NewOutgoingContext(ctx, md)
}

func GetGrpcInContext(ctx context.Context, data map[string]string) context.Context {
	md := metadata.New(data)
	// 发送 metadata
	// 创建带有meta的context
	return metadata.NewIncomingContext(ctx, md)
}
