package live

import (
	"sync"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/conn"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/logger"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	ggrpc "google.golang.org/grpc"
)

const serviceName = "live"

type Client struct {
	lock     sync.RWMutex
	config   config.Config
	registry *etcd.Registry
	conn     *ggrpc.ClientConn
	opts     []ggrpc.DialOption

	liveServiceClient LiveServiceClient
}

func NewClient(cfg config.Config, registry *etcd.Registry, opts ...ggrpc.DialOption) (*Client, func(), error) {
	c := &Client{
		config:   cfg,
		registry: registry,
		opts:     opts,
	}
	if err := c.createConn(); err != nil {
		return nil, nil, err
	}
	cleanFunc := func() {
		if c.conn != nil {
			if err := c.conn.Close(); err != nil {
				logger.Errorf("grpc close error: %s", err.Error())
			}
		}
	}
	return c, cleanFunc, nil
}

func (c *Client) createConn() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.conn != nil {
		return nil
	}
	logger.Infof("%s grpc client conn, %+v", serviceName, c.config)
	cc, err := conn.CreateConn(serviceName, c.config, c.registry, c.opts...)
	if err != nil {
		return errors.NewMsg("grpc.Dial error: %s", err)
	}
	c.liveServiceClient = NewLiveServiceClient(cc)
	c.conn = cc
	return nil
}

func (c *Client) GetLiveServiceClient() (LiveServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.liveServiceClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.liveServiceClient, nil
}
