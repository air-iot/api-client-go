package spm

import (
	"sync"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/conn"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/logger"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"
)

const serviceName = "spm"

type Client struct {
	lock       sync.RWMutex
	conn       *grpc.ClientConn
	restClient *http.Client

	registry      *etcd.Registry
	config        config.Config
	projectClient ProjectServiceClient
	userClient    UserServiceClient
	settingClient SettingServiceClient
	opts          []grpc.DialOption
	middlewares   []middleware.Middleware
}

func NewClient(cfg config.Config, registry *etcd.Registry, cred grpc.DialOption, httpCred middleware.Middleware) (*Client, func(), error) {
	c := &Client{
		registry:    registry,
		config:      cfg,
		opts:        []grpc.DialOption{cred},
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

func (c *Client) createConn() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.conn != nil {
		return nil
	}
	logger.Infof("%s grpc client cc, %+v", serviceName, c.config)
	cc, err := conn.CreateConn(serviceName, c.config, c.registry, c.opts...)
	if err != nil {
		return errors.NewMsg("grpc.Dial error: %s", err)
	}
	c.projectClient = NewProjectServiceClient(cc)
	c.userClient = NewUserServiceClient(cc)
	c.settingClient = NewSettingServiceClient(cc)
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
		return errors.NewMsg("rest error: %s", err)
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

func (c *Client) GetProjectServiceClient() (ProjectServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.projectClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.projectClient, nil
}

func (c *Client) GetUserServiceClient() (UserServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.userClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.userClient, nil
}

func (c *Client) GetSettingServiceClient() (SettingServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.settingClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.settingClient, nil
}
