package sync

import (
	"sync"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/conn"
	"github.com/air-iot/errors"
	"github.com/air-iot/logger"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
	ggrpc "google.golang.org/grpc"
)

const serviceName = "sync"

type Client struct {
	lock sync.RWMutex

	cc ggrpc.ClientConnInterface

	config      config.Config
	registry    *etcd.Registry
	conn        *ggrpc.ClientConn
	restClient  *http.Client
	opts        []ggrpc.DialOption
	middlewares []middleware.Middleware

	syncClient SyncServiceClient
}

func NewClient(cfg config.Config, registry *etcd.Registry, cred ggrpc.DialOption, httpCred middleware.Middleware) (*Client, func(), error) {
	c := &Client{
		registry:    registry,
		config:      cfg,
		opts:        []ggrpc.DialOption{cred},
		middlewares: []middleware.Middleware{httpCred},
	}
	//if err := c.createRestConn(); err != nil {
	//	return nil, nil, err
	//}
	//if err := c.createConn(); err != nil {
	//	return nil, nil, err
	//}
	cleanFunc := func() {
		if c.restClient != nil {
			if err := c.restClient.Close(); err != nil {
				logger.Errorf("rest close error: %s", err.Error())
			}
		}
		if c.conn != nil {
			if err := c.conn.Close(); err != nil {
				logger.Errorf("grpc close error: %s", err.Error())
			}
		}
	}
	return c, cleanFunc, nil
}

func NewLocalClient(cfg config.Config, cc ggrpc.ClientConnInterface) (*Client, func(), error) {
	c := &Client{
		config: cfg,
		cc:     cc,
	}
	c.syncClient = NewSyncServiceClient(cc)
	cleanFunc := func() {}
	return c, cleanFunc, nil
}

func (c *Client) createConn() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.conn != nil {
		return nil
	}
	logger.Infof("%s grpc client cc, %+v", serviceName, c.config)
	cc, err := conn.CreateConn(serviceName, c.config, c.registry, c.opts...)
	if err != nil {
		return err
	}
	c.syncClient = NewSyncServiceClient(cc)
	c.conn = cc
	return nil
}

func (c *Client) createRestConn() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.restClient != nil {
		return nil
	}
	logger.Infof("%s http client createConn, %+v", serviceName, c.config)
	cc, err := conn.CreateRestConn(serviceName, c.config, c.registry, c.middlewares...)
	if err != nil {
		return err
	}
	c.restClient = cc
	return nil
}

func (c *Client) GetRestClient() (*http.Client, error) {
	if c.restClient == nil {
		if err := c.createRestConn(); err != nil {
			return nil, err
		}
	}
	return c.restClient, nil
}

func (c *Client) GetSyncServiceClient() (SyncServiceClient, error) {
	if c.syncClient != nil {
		return c.syncClient, nil
	}
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.syncClient == nil {
		return nil, errors.New("客户端是空")
	}
	return c.syncClient, nil
}
