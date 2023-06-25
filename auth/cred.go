package auth

import (
	"context"
	"fmt"

	"github.com/air-iot/api-client-go/v4/apitransport"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (

	// OpenTLS 是否开启TLS认证
	OpenTLS = false
)

// CustomCredential 自定义认证
type CustomCredential struct {
	f GetAuthClient
}

type GetAuthClient func() *Client

func NewCustomCredential(f GetAuthClient) *CustomCredential {
	return &CustomCredential{f: f}
}

func (c *CustomCredential) HttpToken() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			transporter, ok := transport.FromClientContext(ctx)
			if !ok {
				return nil, errors.NewMsg("客户端上下文错误")
			}
			tt, ok := apitransport.FromClientContext(ctx)
			if ok {
				for k, v := range tt.ReqHeader {
					transporter.RequestHeader().Set(k, v)
				}
			}
			transporter.RequestHeader().Set("Request-Type", "service")
			token := transporter.RequestHeader().Get(config.XRequestHeaderAuthorization)
			if token == "" {
				token, err = (c.f)().Token()
				if err != nil {
					return nil, err
				}
				transporter.RequestHeader().Set(config.XRequestHeaderAuthorization, fmt.Sprintf("Bearer %s", token))
			}
			return handler(ctx, req)
		}
	}
}

// GetRequestMetadata 实现自定义认证接口
func (c *CustomCredential) GetRequestMetadata(ctx context.Context, uri ...string) (mds map[string]string, err error) {
	info, ok := credentials.RequestInfoFromContext(ctx)
	//_, _ = info, ok
	//pr, ok := transport.FromClientContext(ctx)
	if !ok {
		return nil, errors.NewMsg("客户端上下文错误")
	}
	//path := pr.Operation()
	path := info.Method
	if path == "/spm.UserService/GetToken" ||
		path == "/core.LicenseService/FindMachineCode" ||
		path == "/core.AppService/GetToken" {
		//path == "/core.UserService/GetCurrentUserInfo" ||
		//path == "/core.RoleService/AdminRoleCheck" ||
		//path == "/core.TableDataService/GetWarningFilterIDs" ||
		//path == "/warning.WarnService/Query" {
		return map[string]string{}, nil
	}
	md, ok := metadata.FromOutgoingContext(ctx)
	var token string
	if ok {
		headers := md.Get(config.XRequestHeaderAuthorization)
		if len(headers) > 0 {
			if headers[0] != "" {
				return map[string]string{}, nil
			}
		}
	}
	token, err = (c.f)().Token()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		config.XRequestHeaderAuthorization: token,
	}, nil
}

// RequireTransportSecurity 自定义认证是否开启TLS
func (c *CustomCredential) RequireTransportSecurity() bool {
	return OpenTLS
}
