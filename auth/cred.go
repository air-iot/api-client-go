package auth

import (
	"context"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc/metadata"
)

const (

	// OpenTLS 是否开启TLS认证
	OpenTLS = false
)

// customCredential 自定义认证
type customCredential struct {
	f GetAuthClient
}

type GetAuthClient func() *Client

func NewCustomCredential(f GetAuthClient) *customCredential {
	return &customCredential{f: f}
}

// GetRequestMetadata 实现自定义认证接口
func (c *customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (mds map[string]string, err error) {
	pr, ok := transport.FromClientContext(ctx)
	if !ok {
		return nil, errors.NewMsg("客户端上下文错误")
	}
	path := pr.Operation()
	if path == "/spm.UserService/GetToken" ||
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
			return map[string]string{}, nil
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
func (c *customCredential) RequireTransportSecurity() bool {
	return OpenTLS
}
