package eap

import (
	"sync"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/conn"
	"github.com/air-iot/errors"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"google.golang.org/grpc"
)

const serviceName = "eap"

type Client struct {
	lock       sync.RWMutex
	cc         grpc.ClientConnInterface
	conn       *grpc.ClientConn
	registry   *etcd.Registry
	config     config.Config
	opts       []grpc.DialOption
	taskClient TaskServiceClient
}

func NewClient(cfg config.Config, registry *etcd.Registry, cred grpc.DialOption) (*Client, func(), error) {
	c := &Client{registry: registry, config: cfg, opts: []grpc.DialOption{cred}}
	return c, func() {
		if c.conn != nil {
			_ = c.conn.Close()
		}
	}, nil
}

func NewLocalClient(cfg config.Config, cc grpc.ClientConnInterface) (*Client, func(), error) {
	return &Client{config: cfg, cc: cc, taskClient: NewTaskServiceClient(cc)}, func() {}, nil
}

func (c *Client) createConn() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.conn != nil {
		return nil
	}
	cc, err := conn.CreateConn(serviceName, c.config, c.registry, c.opts...)
	if err != nil {
		return err
	}
	c.conn = cc
	c.taskClient = NewTaskServiceClient(cc)
	return nil
}

func (c *Client) GetTaskServiceClient() (TaskServiceClient, error) {
	if c.taskClient != nil {
		return c.taskClient, nil
	}
	if err := c.createConn(); err != nil {
		return nil, err
	}
	if c.taskClient == nil {
		return nil, errors.New("客户端是空")
	}
	return c.taskClient, nil
}
