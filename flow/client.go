package flow

import (
	"fmt"
	conn2 "github.com/air-iot/api-client-go/v4/conn"
	"sync"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	ggrpc "google.golang.org/grpc"

	"github.com/air-iot/logger"

	"github.com/air-iot/api-client-go/v4/config"
)

const serviceName = "flow"

type Client struct {
	lock     sync.RWMutex
	config   config.Config
	registry *etcd.Registry
	conn     *ggrpc.ClientConn

	flowTaskClient                 FlowTaskServiceClient
	flowClient                     FlowServiceClient
	flowTriggerRecordServiceClient FlowTriggerRecordServiceClient
}

func NewClient(cfg config.Config, registry *etcd.Registry) (*Client, func(), error) {
	c := &Client{
		config:   cfg,
		registry: registry,
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
	logger.Infof("flow grpc client conn, %+v", c.config)
	conn, err := conn2.CreateConn(serviceName, c.config, c.registry)
	if err != nil {
		return fmt.Errorf("grpc.Dial error: %s", err)
	}
	if err != nil {
		return fmt.Errorf("grpc.Dial error: %s", err)
	}
	c.conn = conn
	c.flowTaskClient = NewFlowTaskServiceClient(conn)
	c.flowClient = NewFlowServiceClient(conn)
	c.flowTriggerRecordServiceClient = NewFlowTriggerRecordServiceClient(conn)
	return nil
}

func (c *Client) GetFlowServiceClient() (FlowServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.flowClient == nil {
		return nil, fmt.Errorf("客户端是空")
	}
	return c.flowClient, nil
}

func (c *Client) GetFlowTaskServiceClient() (FlowTaskServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.flowTaskClient == nil {
		return nil, fmt.Errorf("客户端是空")
	}
	return c.flowTaskClient, nil
}

func (c *Client) GetFlowTriggerRecordServiceClient() (FlowTriggerRecordServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.flowTriggerRecordServiceClient == nil {
		return nil, fmt.Errorf("客户端是空")
	}
	return c.flowTriggerRecordServiceClient, nil
}
