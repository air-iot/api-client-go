package conn

import (
	"context"
	"fmt"
	"io"
	netHttp "net/http"
	"time"

	"github.com/air-iot/logger"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	ggrpc "google.golang.org/grpc"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/api-client-go/v4/filter"
)

func CreateConn(serviceName string, cfg config.Config, r *etcd.Registry, opts ...ggrpc.DialOption) (*ggrpc.ClientConn, error) {
	metadataTmp := cfg.Metadata
	if srv, ok := cfg.Services[serviceName]; ok {
		if srv.Metadata != nil && len(srv.Metadata) > 0 {
			metadataTmp = srv.Metadata
		}
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 60
	}
	if cfg.Limit == 0 {
		cfg.Limit = 100
	}
	opts = append(opts, ggrpc.WithDefaultCallOptions(ggrpc.MaxCallRecvMsgSize(cfg.Limit*1024*1024), ggrpc.MaxCallSendMsgSize(cfg.Limit*1024*1024)))
	logger.Infof("grpc create conn,serviceName: %s config: %+v", serviceName, cfg)
	return grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", serviceName)),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
		grpc.WithOptions(opts...),
		grpc.WithPrintDiscoveryDebugLog(cfg.Debug),
		grpc.WithNodeFilter(filter.Metadata(metadataTmp)),
		grpc.WithTimeout(time.Second*time.Duration(cfg.Timeout)),
	)
}

func CreateRestConn(serviceName string, cfg config.Config, r *etcd.Registry, middlewares ...middleware.Middleware) (*http.Client, error) {
	metadataTmp := cfg.Metadata
	if srv, ok := cfg.Services[serviceName]; ok {
		if srv.Metadata != nil && len(srv.Metadata) > 0 {
			metadataTmp = srv.Metadata
		}
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 60
	}
	logger.Infof("http create conn,serviceName: %s config: %+v", serviceName, cfg)
	return http.NewClient(
		context.Background(),
		http.WithEndpoint(fmt.Sprintf("discovery:///%s", serviceName)),
		http.WithDiscovery(r),
		http.WithMiddleware(
			recovery.Recovery(),
		),
		http.WithMiddleware(middlewares...),
		http.WithErrorDecoder(func(ctx context.Context, res *netHttp.Response) error {
			if res.StatusCode >= 200 && res.StatusCode <= 299 {
				return nil
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(res.Body)
			data, err := io.ReadAll(res.Body)
			if err == nil {
				return errors.ParseBody(res.StatusCode, data)
			}
			return errors.NewMsg("未知原因,解析响应错误,%v", err)
		}),
		http.WithNodeFilter(filter.Metadata(metadataTmp)),
		http.WithTimeout(time.Second*time.Duration(cfg.Timeout)),
	)
}
