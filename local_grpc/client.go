package local_grpc

import (
	"context"
	"github.com/air-iot/api-client-go/v4/auth"
	"github.com/air-iot/api-client-go/v4/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Client struct {
	s *Server
	f auth.GetAuthClient
}

func NewClient(s *Server, f auth.GetAuthClient) grpc.ClientConnInterface {
	return &Client{s: s, f: f}
}

func (cc *Client) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	metadataIn, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		headers := metadataIn.Get(config.XRequestHeaderAuthorization)
		if len(headers) == 0 && cc.f != nil {
			token, err := (cc.f)().Token()
			if err != nil {
				return err
			}
			metadataIn.Set(config.XRequestHeaderAuthorization, token)
		}
		ctx = metadata.NewIncomingContext(ctx, metadataIn)
	}
	return cc.s.Invoke(ctx, method, args, reply)
}

func (cc *Client) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
