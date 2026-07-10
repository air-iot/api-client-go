package conn

import (
	"context"
	"fmt"
	"io"
	netHttp "net/http"
	"time"

	"github.com/air-iot/errors"
	"github.com/air-iot/logger"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	ggrpc "google.golang.org/grpc"

	"github.com/air-iot/api-client-go/v4/config"
	internalError "github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/api-client-go/v4/filter"
)

func CreateConn(serviceName string, cfg config.Config, r *etcd.Registry, opts ...ggrpc.DialOption) (*ggrpc.ClientConn, error) {
	metadataTmp := cfg.Metadata
	addr := cfg.GatewayGrpc
	if srv, ok := cfg.Services[serviceName]; ok {
		if srv.Metadata != nil && len(srv.Metadata) > 0 {
			metadataTmp = srv.Metadata
		}
		if srv.GrpcAddr != "" {
			addr = srv.GrpcAddr
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
	if r != nil {
		// WithBlock: 阻塞到 discovery 首次解析出可用节点后再返回,
		// 避免 gRPC 客户端懒创建时,首次请求早于服务发现就绪而报 no_available_node。
		opts = append(opts, ggrpc.WithBlock())
		cli, err := grpc.DialInsecure(
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
		if err != nil {
			return nil, errors.Wrap(err, "grpc.Dial err")
		}
		return cli, nil
	} else {

		cli, err := grpc.DialInsecure(
			context.Background(),
			grpc.WithEndpoint(addr),
			grpc.WithMiddleware(
				tracing.Client(),
				recovery.Recovery(),
			),
			grpc.WithOptions(opts...),
			grpc.WithPrintDiscoveryDebugLog(cfg.Debug),
			//grpc.WithNodeFilter(filter.Metadata(metadataTmp)),
			grpc.WithTimeout(time.Second*time.Duration(cfg.Timeout)),
		)
		if err != nil {
			return nil, errors.Wrap(err, "grpc.Dial err")
		}
		return cli, nil
	}
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
	opts := []http.ClientOption{
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
				return internalError.ParseBody(res.StatusCode, data)
			}
			return errors.Wrap(err, "未知原因,解析响应错误")
		}),
		http.WithNodeFilter(filter.Metadata(metadataTmp)),
		http.WithTimeout(time.Second * time.Duration(cfg.Timeout)),
		//http.WithTransport(tr),
	}
	if r != nil {
		opts = append(opts,
			http.WithEndpoint(fmt.Sprintf("discovery:///%s", serviceName)),
			http.WithDiscovery(r),
			// WithBlock: 阻塞到 discovery 首次解析出可用节点后再返回,
			// 避免 REST 客户端懒创建时,首次请求早于服务发现就绪而报 no_available_node。
			http.WithBlock(),
		)
		//cli, err := http.NewClient(
		//	context.Background(),
		//	http.WithEndpoint(fmt.Sprintf("discovery:///%s", serviceName)),
		//	http.WithDiscovery(r),
		//	http.WithMiddleware(
		//		recovery.Recovery(),
		//	),
		//	http.WithMiddleware(middlewares...),
		//	http.WithErrorDecoder(func(ctx context.Context, res *netHttp.Response) error {
		//		if res.StatusCode >= 200 && res.StatusCode <= 299 {
		//			return nil
		//		}
		//		defer func(Body io.ReadCloser) {
		//			err := Body.Close()
		//			if err != nil {
		//
		//			}
		//		}(res.Body)
		//		data, err := io.ReadAll(res.Body)
		//		if err == nil {
		//			return internalError.ParseBody(res.StatusCode, data)
		//		}
		//		return errors.Wrap(err, "未知原因,解析响应错误")
		//	}),
		//	http.WithNodeFilter(filter.Metadata(metadataTmp)),
		//	http.WithTimeout(time.Second*time.Duration(cfg.Timeout)),
		//	http.WithTransport(tr),
		//)
		//
		//if err != nil {
		//	return nil, errors.Wrap(err, "create http client err")
		//}
		//return cli, nil
	} else {
		opts = append(opts, http.WithEndpoint(cfg.Gateway))
		//cli, err := http.NewClient(
		//	context.Background(),
		//	http.WithEndpoint(cfg.Gateway),
		//	http.WithMiddleware(
		//		recovery.Recovery(),
		//	),
		//	http.WithMiddleware(middlewares...),
		//	http.WithErrorDecoder(func(ctx context.Context, res *netHttp.Response) error {
		//		if res.StatusCode >= 200 && res.StatusCode <= 299 {
		//			return nil
		//		}
		//		defer func(Body io.ReadCloser) {
		//			err := Body.Close()
		//			if err != nil {
		//
		//			}
		//		}(res.Body)
		//		data, err := io.ReadAll(res.Body)
		//		if err == nil {
		//			return internalError.ParseBody(res.StatusCode, data)
		//		}
		//		return errors.Wrap(err, "未知原因,解析响应错误")
		//	}),
		//	http.WithNodeFilter(filter.Metadata(metadataTmp)),
		//	http.WithTimeout(time.Second*time.Duration(cfg.Timeout)),
		//	http.WithTransport(tr),
		//)

		//if err != nil {
		//	return nil, errors.Wrap(err, "create http client err")
		//}
		//return cli, nil
	}

	if cfg.KeepAlive {
		tr := &netHttp.Transport{
			DisableKeepAlives: !cfg.KeepAlive, // Explicitly enable keep-alives (default for HTTP/1.1)
			MaxIdleConns:      cfg.MaxIdleConns,
			IdleConnTimeout:   cfg.IdleConnTimeout,
			// Other transport configurations can go here, e.g., TLSClientConfig, Proxy
		}
		opts = append(opts, http.WithTransport(tr))
	}

	cli, err := http.NewClient(
		context.Background(),
		opts...,
	)

	if err != nil {
		return nil, errors.Wrap(err, "create http client err")
	}
	return cli, nil
}
