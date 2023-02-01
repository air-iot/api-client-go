package conn

import (
	"context"
	"fmt"
	"time"

	"github.com/air-iot/logger"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/filter"
)

// Config grpc配置参数
//type Config struct {
//	Host string
//	Port int
//	AK   string
//	SK   string
//	//Timeout int
//}

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
	logger.Infof("grpc create conn,serviceName: %s config: %+v", serviceName, cfg)
	return grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", serviceName)),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
		grpc.WithOptions(opts...),
		grpc.WithNodeFilter(filter.Metadata(metadataTmp)),
		grpc.WithTimeout(time.Second*time.Duration(cfg.Timeout)),
	)
}
