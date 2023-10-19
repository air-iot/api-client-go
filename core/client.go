package core

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

const serviceName = "core"

type Client struct {
	lock        sync.RWMutex
	conn        *grpc.ClientConn
	restClient  *http.Client
	config      config.Config
	registry    *etcd.Registry
	opts        []grpc.DialOption
	middlewares []middleware.Middleware

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
	backupServiceClient        BackupServiceClient
	taskManagerServiceClient   TaskManagerServiceClient
	dashboardClient            DashboardServiceClient
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
	c.backupServiceClient = NewBackupServiceClient(cc)
	c.dashboardClient = NewDashboardServiceClient(cc)
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

func (c *Client) GetBackupServiceClient() (BackupServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.backupServiceClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.backupServiceClient, nil
}

func (c *Client) GetTaskManagerServiceClient() (TaskManagerServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.taskManagerServiceClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.taskManagerServiceClient, nil
}

func (c *Client) GetDashboardServiceClient() (DashboardServiceClient, error) {
	if c.conn == nil {
		if err := c.createConn(); err != nil {
			return nil, err
		}
	}
	if c.dashboardClient == nil {
		return nil, errors.NewMsg("客户端是空")
	}
	return c.dashboardClient, nil
}
