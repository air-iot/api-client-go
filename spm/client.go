package spm

import (
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"google.golang.org/grpc"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/conn"
	"github.com/air-iot/logger"
)

const serviceName = "spm"

type Client struct {
	lock          sync.RWMutex
	conn          *grpc.ClientConn
	registry      *etcd.Registry
	config        config.Config
	projectClient ProjectServiceClient
}

func NewClient(cfg config.Config, registry *etcd.Registry) (*Client, func(), error) {
	c := &Client{
		registry: registry,
		config:   cfg,
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
	logger.Infof("flow grpc client createConn, %+v", c.config)
	createConn, err := conn.CreateConn(serviceName, c.config, c.registry)
	if err != nil {
		return fmt.Errorf("grpc.Dial error: %s", err)
	}
	if err != nil {
		return fmt.Errorf("grpc.Dial error: %s", err)
	}
	c.conn = createConn
	c.projectClient = NewProjectServiceClient(createConn)
	return nil
}

func (c *Client) GetProjectServiceClient() (ProjectServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.projectClient == nil {
		return nil, fmt.Errorf("客户端是空")
	}
	return c.projectClient, nil
}
