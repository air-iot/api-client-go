package report

import (
	"sync"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/conn"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/logger"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"google.golang.org/grpc"
)

const serviceName = "report"

type Client struct {
	lock     sync.RWMutex
	config   config.Config
	registry *etcd.Registry
	conn     *grpc.ClientConn
	opts     []grpc.DialOption

	reportClient ReportServiceClient
	reportCopyClient ReportCopyServiceClient
}

func NewClient(cfg config.Config, registry *etcd.Registry, opts ...grpc.DialOption) (*Client, func(), error) {
	c := &Client{
		registry: registry,
		config:   cfg,
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
	logger.Infof("flow grpc client cc, %+v", c.config)
	cc, err := conn.CreateConn(serviceName, c.config, c.registry, c.opts...)
	if err != nil {
		return errors.NewMsg("grpc.Dial error: %s", err)
	}
	c.reportClient = NewReportServiceClient(cc)
	c.reportCopyClient = NewReportCopyServiceClient(cc)
	c.conn = cc
	return nil
}

func (c *Client) GetReportServiceClient() (ReportServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.reportClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.reportClient, nil
}

func (c *Client) GetReportCopyServiceClient() (ReportCopyServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.reportCopyClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.reportCopyClient, nil
}
