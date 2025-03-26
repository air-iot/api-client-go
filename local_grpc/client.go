package local_grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/air-iot/api-client-go/v4/auth"
	"github.com/air-iot/api-client-go/v4/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
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
	sm := method
	if sm != "" && sm[0] == '/' {
		sm = sm[1:]
	}
	pos := strings.LastIndex(sm, "/")
	if pos == -1 {
		errDesc := fmt.Sprintf("malformed method name: %q", method)
		return status.Error(codes.Unimplemented, errDesc)
	}
	service := sm[:pos]
	method1 := sm[pos+1:]
	srv, knownService := cc.s.GetService(service)
	codec := encoding.GetCodec(proto.Name)
	if knownService {
		if md, ok := srv.GetMethods(method1); ok {
			df := func(v any) error {
				d, err := codec.Marshal(args)
				if err != nil {
					return status.Errorf(codes.Internal, "grpc: error marshalling request: %v", err)
				}
				if err := codec.Unmarshal(d, v); err != nil {
					return status.Errorf(codes.Internal, "grpc: error unmarshalling request: %v", err)
				}
				return nil
			}
			mdIn, ok := metadata.FromOutgoingContext(ctx)
			if ok {
				headers := mdIn.Get(config.XRequestHeaderAuthorization)
				if len(headers) == 0 && cc.f != nil {
					token, err := (cc.f)().Token()
					if err != nil {
						return err
					}
					mdIn.Set(config.XRequestHeaderAuthorization, token)
				}
				ctx = metadata.NewIncomingContext(ctx, mdIn)
			}
			appReply, appErr := md.Handler(srv.GetServiceImpl(), ctx, df, nil)
			if appErr != nil {
				appStatus, ok := status.FromError(appErr)
				if !ok {
					// Convert non-status application error to a status error with code
					// Unknown, but handle context errors specifically.
					appStatus = status.FromContextError(appErr)
					appErr = appStatus.Err()
				}
				return appErr
			}

			b, err := codec.Marshal(appReply)
			if err != nil {
				return status.Errorf(codes.Internal, "grpc: error while marshaling: %v", err.Error())
			}

			if err := codec.Unmarshal(b, reply); err != nil {
				return status.Errorf(codes.Internal, "grpc: error while unmarshaling: %v", err.Error())
			}
			return nil
		}
	}
	errDesc := fmt.Sprintf("unknown service %v", service)
	return status.Error(codes.Unimplemented, errDesc)
}

func (cc *Client) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
