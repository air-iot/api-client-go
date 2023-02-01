package core

import (
	"sync"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/conn"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/logger"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"google.golang.org/grpc"
)

const serviceName = "core"

type Client struct {
	lock sync.RWMutex

	conn     *grpc.ClientConn
	config   config.Config
	registry *etcd.Registry

	opts []grpc.DialOption

	appServiceClient           AppServiceClient
	licenseServiceClient       LicenseServiceClient
	userServiceClient          UserServiceClient
	logServiceClient           LogServiceClient
	tableSchemaClient          TableSchemaServiceClient
	tableRecordClient          TableRecordServiceClient
	tableDataClient            TableDataServiceClient
	messageClient              MessageServiceClient
	dataQueryClient            DataQueryServiceClient
	roleClient                 RoleServiceClient
	catalogClient              CatalogServiceClient
	deptClient                 DeptServiceClient
	settingClient              SettingServiceClient
	systemVariablServiceClient SystemVariableServiceClient
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
	logger.Infof("core grpc client cc, %+v", c.config)
	cc, err := conn.CreateConn(serviceName, c.config, c.registry, c.opts...)
	if err != nil {
		return errors.NewMsg("grpc.Dial error: %s", err)
	}
	c.appServiceClient = NewAppServiceClient(cc)
	c.licenseServiceClient = NewLicenseServiceClient(cc)
	c.logServiceClient = NewLogServiceClient(cc)
	c.userServiceClient = NewUserServiceClient(cc)
	c.tableSchemaClient = NewTableSchemaServiceClient(cc)
	c.tableRecordClient = NewTableRecordServiceClient(cc)
	c.tableDataClient = NewTableDataServiceClient(cc)
	c.messageClient = NewMessageServiceClient(cc)
	c.dataQueryClient = NewDataQueryServiceClient(cc)
	c.roleClient = NewRoleServiceClient(cc)
	c.catalogClient = NewCatalogServiceClient(cc)
	c.deptClient = NewDeptServiceClient(cc)
	c.settingClient = NewSettingServiceClient(cc)
	c.systemVariablServiceClient = NewSystemVariableServiceClient(cc)
	c.conn = cc
	return nil
}

func (c *Client) GetAppServiceClient() (AppServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.appServiceClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.appServiceClient, nil
}

func (c *Client) GetLicenseServiceClient() (LicenseServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.licenseServiceClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.licenseServiceClient, nil
}

func (c *Client) GetLogServiceClient() (LogServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.logServiceClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.logServiceClient, nil
}

func (c *Client) GetUserServiceClient() (UserServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.userServiceClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.userServiceClient, nil
}

func (c *Client) GetTableSchemaServiceClient() (TableSchemaServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.tableSchemaClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.tableSchemaClient, nil
}

func (c *Client) GetTableRecordServiceClient() (TableRecordServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.tableRecordClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.tableRecordClient, nil
}

func (c *Client) GetTableDataServiceClient() (TableDataServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.tableDataClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.tableDataClient, nil
}

func (c *Client) GetMessageServiceClient() (MessageServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.messageClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.messageClient, nil
}

func (c *Client) GetDataQueryServiceClient() (DataQueryServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.dataQueryClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.dataQueryClient, nil
}

func (c *Client) GetRoleServiceClient() (RoleServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.roleClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.roleClient, nil
}

func (c *Client) GetCatalogServiceClient() (CatalogServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.catalogClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.catalogClient, nil
}

func (c *Client) GetDeptServiceClient() (DeptServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.deptClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.deptClient, nil
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

func (c *Client) GetSystemVariableServiceClient() (SystemVariableServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.systemVariablServiceClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.systemVariablServiceClient, nil
}
