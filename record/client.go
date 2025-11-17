package record

import (
	"context"
	"fmt"
	"sync"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/apicontext"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/conn"
	"github.com/air-iot/errors"
	"github.com/air-iot/logger"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
	ggrpc "google.golang.org/grpc"
)

const serviceName = "record"

type Client struct {
	lock sync.RWMutex

	cc ggrpc.ClientConnInterface

	config      config.Config
	registry    *etcd.Registry
	conn        *ggrpc.ClientConn
	restClient  *http.Client
	opts        []ggrpc.DialOption
	middlewares []middleware.Middleware

	recordClient RecordServiceClient
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
	c.recordClient = NewRecordServiceClient(cc)
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
	c.recordClient = NewRecordServiceClient(cc)
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

func (c *Client) GetRecordServiceClient() (RecordServiceClient, error) {
	if c.recordClient != nil {
		return c.recordClient, nil
	}
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.recordClient == nil {
		return nil, errors.New("客户端是空")
	}
	return c.recordClient, nil
}

// DeleteWarning 删除报警录制相关信息. 包括录制视频和抓图信息
func (c *Client) DeleteWarning(ctx context.Context, projectId, warningId string, deleteFiles bool, opts ...ggrpc.CallOption) (*api.Response, error) {
	cli, err := c.GetRecordServiceClient()
	if err != nil {
		return nil, err
	}

	if projectId == "" {
		return nil, fmt.Errorf("项目ID不能为空")
	} else if warningId == "" {
		return nil, fmt.Errorf("报警ID不能为空")
	}

	return cli.DeleteWarning(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&DeleteWarningRecordRequest{
			WarningId:   warningId,
			DeleteFiles: deleteFiles,
		}, opts...)
}

// BatchDeleteWarnings 批量删除报警录制相关信息. 包括录制视频和抓图信息
func (c *Client) BatchDeleteWarnings(ctx context.Context, projectId string, warningIds []string, deleteFiles bool, opts ...ggrpc.CallOption) (*api.Response, error) {
	cli, err := c.GetRecordServiceClient()
	if err != nil {
		return nil, err
	}

	if projectId == "" {
		return nil, fmt.Errorf("项目ID不能为空")
	} else if len(warningIds) == 0 {
		return nil, fmt.Errorf("报警ID不能为空")
	}

	return cli.BatchDeleteWarnings(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&BatchDeleteWarningRecordRequest{
			WarningIds:  warningIds,
			DeleteFiles: deleteFiles,
		}, opts...)
}

// DeleteWarningPlayback 删除报警录制视频信息
func (c *Client) DeleteWarningPlayback(ctx context.Context, projectId, warningId string, deleteFiles bool, opts ...ggrpc.CallOption) (*api.Response, error) {
	cli, err := c.GetRecordServiceClient()
	if err != nil {
		return nil, err
	}

	if projectId == "" {
		return nil, fmt.Errorf("项目ID不能为空")
	} else if warningId == "" {
		return nil, fmt.Errorf("报警ID不能为空")
	}

	return cli.DeleteWarningPlayback(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &DeleteWarningRecordRequest{
		WarningId:   warningId,
		DeleteFiles: deleteFiles,
	}, opts...)
}

// BatchDeleteWarningPlaybacks 批量删除报警录制视频信息
func (c *Client) BatchDeleteWarningPlaybacks(ctx context.Context, projectId string, warningIds []string, deleteFiles bool, opts ...ggrpc.CallOption) (*api.Response, error) {
	cli, err := c.GetRecordServiceClient()
	if err != nil {
		return nil, err
	}

	if projectId == "" {
		return nil, fmt.Errorf("项目ID不能为空")
	} else if len(warningIds) == 0 {
		return nil, fmt.Errorf("报警ID不能为空")
	}

	return cli.BatchDeleteWarningPlaybacks(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &BatchDeleteWarningRecordRequest{
		WarningIds:  warningIds,
		DeleteFiles: deleteFiles,
	}, opts...)
}

// DeleteWarningCapture 删除报警录制抓图信息
func (c *Client) DeleteWarningCapture(ctx context.Context, projectId, warningId string, deleteFiles bool, opts ...ggrpc.CallOption) (*api.Response, error) {
	cli, err := c.GetRecordServiceClient()
	if err != nil {
		return nil, err
	}

	if projectId == "" {
		return nil, fmt.Errorf("项目ID不能为空")
	} else if warningId == "" {
		return nil, fmt.Errorf("报警ID不能为空")
	}

	return cli.DeleteWarningCapture(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &DeleteWarningRecordRequest{
		WarningId:   warningId,
		DeleteFiles: deleteFiles,
	}, opts...)
}

// BatchDeleteWarningCaptures 批量删除报警录制抓图信息
func (c *Client) BatchDeleteWarningCaptures(ctx context.Context, projectId string, warningIds []string, deleteFiles bool, opts ...ggrpc.CallOption) (*api.Response, error) {
	cli, err := c.GetRecordServiceClient()
	if err != nil {
		return nil, err
	}

	if projectId == "" {
		return nil, fmt.Errorf("项目ID不能为空")
	} else if len(warningIds) == 0 {
		return nil, fmt.Errorf("报警ID不能为空")
	}

	return cli.BatchDeleteWarningCaptures(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &BatchDeleteWarningRecordRequest{
		WarningIds:  warningIds,
		DeleteFiles: deleteFiles,
	}, opts...)
}
