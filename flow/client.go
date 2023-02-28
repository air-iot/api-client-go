package flow

import (
	"sync"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/conn"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/logger"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	ggrpc "google.golang.org/grpc"
)

const serviceName = "flow"

type Client struct {
	lock     sync.RWMutex
	config   config.Config
	registry *etcd.Registry
	conn     *ggrpc.ClientConn
	opts     []ggrpc.DialOption

	flowTaskClient                 FlowTaskServiceClient
	flowClient                     FlowServiceClient
	flowJobClient                  FlowJobServiceClient
	flowTriggerRecordServiceClient FlowTriggerRecordServiceClient
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
	logger.Infof("flow grpc client conn, %+v", c.config)
	cc, err := conn.CreateConn(serviceName, c.config, c.registry, c.opts...)
	if err != nil {
		return errors.NewMsg("grpc.Dial error: %s", err)
	}
	c.flowTaskClient = NewFlowTaskServiceClient(cc)
	c.flowClient = NewFlowServiceClient(cc)
	c.flowJobClient = NewFlowJobServiceClient(cc)
	c.flowTriggerRecordServiceClient = NewFlowTriggerRecordServiceClient(cc)
	c.conn = cc
	return nil
}

func (c *Client) GetFlowServiceClient() (FlowServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.flowClient == nil {
		return nil, errors.NewMsg("客户端是空")
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
		return nil, errors.NewMsg("客户端是空")
	}
	return c.flowTaskClient, nil
}

func (c *Client) GetFlowJobServiceClient() (FlowJobServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.flowJobClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.flowJobClient, nil
}

func (c *Client) GetFlowTriggerRecordServiceClient() (FlowTriggerRecordServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.flowTriggerRecordServiceClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.flowTriggerRecordServiceClient, nil
}
