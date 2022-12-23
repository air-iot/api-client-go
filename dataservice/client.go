package dataservice

import (
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"google.golang.org/grpc"

	"github.com/air-iot/logger"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/conn"
)

const serviceName = "data-service"

type Client struct {
	lock     sync.RWMutex
	config   config.Config
	registry *etcd.Registry
	conn     *grpc.ClientConn

	dataServiceClient DataServiceClient
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
	logger.Infof("data-service grpc client conn, %+v", c.config)
	cc, err := conn.CreateConn(serviceName, c.config, c.registry)
	if err != nil {
		return fmt.Errorf("grpc.Dial error: %s", err)
	}
	if err != nil {
		return fmt.Errorf("grpc.Dial error: %s", err)
	}
	c.conn = cc
	c.dataServiceClient = NewDataServiceClient(cc)
	return nil
}

func (c *Client) GetDataServiceClient() (DataServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.dataServiceClient == nil {
		return nil, fmt.Errorf("客户端是空")
	}
	return c.dataServiceClient, nil
}
