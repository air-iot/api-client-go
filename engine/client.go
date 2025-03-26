package engine

import (
	"sync"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/conn"
	"github.com/air-iot/errors"
	"github.com/air-iot/logger"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/health"
)

//
//type Config struct {
//	Host string
//	Port int
//}

const serviceName = "flow-engine"

type FlowRunResponse struct {
	Job string `json:"job"`
}

type Client struct {
	lock sync.RWMutex

	cc grpc.ClientConnInterface

	config      config.Config
	registry    *etcd.Registry
	conn        *grpc.ClientConn
	restClient  *http.Client
	opts        []grpc.DialOption
	middlewares []middleware.Middleware

	engineServiceClient      EngineServiceClient
	pluginServiceClient      PluginServiceClient
	flowJobCronServiceClient FlowJobCronServiceClient
	flowLogCronServiceClient FlowLogCronServiceClient
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

func NewLocalClient(cfg config.Config, cc grpc.ClientConnInterface) (*Client, func(), error) {
	c := &Client{
		config: cfg,
		cc:     cc,
	}
	c.engineServiceClient = NewEngineServiceClient(cc)
	c.pluginServiceClient = NewPluginServiceClient(cc)
	c.flowJobCronServiceClient = NewFlowJobCronServiceClient(cc)
	c.flowLogCronServiceClient = NewFlowLogCronServiceClient(cc)
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
	c.engineServiceClient = NewEngineServiceClient(cc)
	c.pluginServiceClient = NewPluginServiceClient(cc)
	c.flowJobCronServiceClient = NewFlowJobCronServiceClient(cc)
	c.flowLogCronServiceClient = NewFlowLogCronServiceClient(cc)
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

func (c *Client) GetDataServiceClient() (EngineServiceClient, error) {
	if c.engineServiceClient != nil {
		return c.engineServiceClient, nil
	}
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.engineServiceClient == nil {
		return nil, errors.New("客户端是空")
	}
	return c.engineServiceClient, nil
}

func (c *Client) GetPluginServiceClient() (PluginServiceClient, error) {
	if c.pluginServiceClient != nil {
		return c.pluginServiceClient, nil
	}
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.pluginServiceClient == nil {
		return nil, errors.New("客户端是空")
	}
	return c.pluginServiceClient, nil
}

func (c *Client) GetFlowJobCronServiceClient() (FlowJobCronServiceClient, error) {
	if c.flowJobCronServiceClient != nil {
		return c.flowJobCronServiceClient, nil
	}
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.flowJobCronServiceClient == nil {
		return nil, errors.New("客户端是空")
	}
	return c.flowJobCronServiceClient, nil
}

func (c *Client) GetFlowLogCronServiceClient() (FlowLogCronServiceClient, error) {
	if c.flowLogCronServiceClient != nil {
		return c.flowLogCronServiceClient, nil
	}
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.flowLogCronServiceClient == nil {
		return nil, errors.New("客户端是空")
	}
	return c.flowLogCronServiceClient, nil
}
