package api_client_go

import (
	"context"
	"sync"
	"time"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/core"
	"github.com/air-iot/api-client-go/v4/dataservice"
	"github.com/air-iot/api-client-go/v4/driver"
	"github.com/air-iot/api-client-go/v4/engine"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/api-client-go/v4/flow"
	"github.com/air-iot/api-client-go/v4/metadata"
	"github.com/air-iot/api-client-go/v4/spm"
	"github.com/air-iot/api-client-go/v4/warning"
	"github.com/air-iot/json"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	lock sync.RWMutex

	SpmClient         *spm.Client
	CoreClient        *core.Client
	FlowClient        *flow.Client
	WarningClient     *warning.Client
	DriverClient      *driver.Client
	DataServiceClient *dataservice.Client
	FlowEngineClient  *engine.Client

	Config    config.Config
	authToken *AuthToken
	tokens    sync.Map
}

func NewClient(cli *clientv3.Client, cfg config.Config) (*Client, func(), error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 120
	}
	r := etcd.New(cli)
	spmClient, cleanSpm, err := spm.NewClient(cfg, r)
	if err != nil {
		return nil, nil, err
	}
	coreClient, cleanCore, err := core.NewClient(cfg, r)
	if err != nil {
		return nil, nil, err
	}
	flowClient, cleanFlow, err := flow.NewClient(cfg, r)
	if err != nil {
		return nil, nil, err
	}
	warningClient, cleanWarning, err := warning.NewClient(cfg, r)
	if err != nil {
		return nil, nil, err
	}
	driverClient, cleanDriver, err := driver.NewClient(cfg, r)
	if err != nil {
		return nil, nil, err
	}
	dataServiceClient, cleanDataService, err := dataservice.NewClient(cfg, r)
	if err != nil {
		return nil, nil, err
	}
	flowEngineClient, cleanFlowEngine, err := engine.NewClient(cfg, r)
	return &Client{
			SpmClient:         spmClient,
			CoreClient:        coreClient,
			FlowClient:        flowClient,
			WarningClient:     warningClient,
			DriverClient:      driverClient,
			DataServiceClient: dataServiceClient,
			FlowEngineClient:  flowEngineClient,
			Config:            cfg,
			tokens:            sync.Map{},
		}, func() {
			cleanSpm()
			cleanCore()
			cleanFlow()
			cleanWarning()
			cleanDriver()
			cleanDataService()
			cleanFlowEngine()
		}, nil
}

type AuthToken struct {
	TokenType   string `json:"tokenType"`
	ExpiresAt   int64  `json:"expiresAt"`
	AccessToken string `json:"accessToken"`
}

func (c *Client) GetToken(key string) (a *AuthToken, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var res *api.Response
	if key == string(config.Tenant) {
		cli, err := c.SpmClient.GetUserServiceClient()
		if err != nil {
			return nil, errors.NewMsg("获取客户端错误,%s", err)
		}
		res, err = cli.GetToken(context.Background(), &api.TokenRequest{Ak: c.Config.AK, Sk: c.Config.SK})
	} else {
		cli, err := c.CoreClient.GetAppServiceClient()
		if err != nil {
			return nil, errors.NewMsg("获取客户端错误,%s", err)
		}
		res, err = cli.GetToken(
			metadata.GetGrpcContext(context.Background(), map[string]string{config.XRequestProject: key}),
			&api.TokenRequest{Ak: c.Config.AK, Sk: c.Config.SK})
	}
	if err != nil {
		return nil, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return nil, errors.NewMsg("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	var authToken AuthToken
	if err := json.Unmarshal(res.GetResult(), &authToken); err != nil {
		return nil, errors.NewMsg("解析 token 请求结果错误, %s", err)
	}
	c.tokens.Store(key, &authToken)
	return &authToken, nil
}

func (c *Client) Token(projectId string) (token string, err error) {
	var authToken *AuthToken
	var key string
	switch c.Config.Type {
	case config.Project:
		key = projectId
	case config.Tenant:
		key = string(config.Tenant)
	default:
		return "", errors.NewMsg("未知ak sk类型")
	}
	tokenValue, ok := c.tokens.Load(key)
	if !ok {
		authToken, err = c.GetToken(key)
		if err != nil {
			return "", err
		}
	} else {
		authToken = tokenValue.(*AuthToken)
		if authToken.ExpiresAt <= time.Now().Unix() {
			authToken, err = c.GetToken(key)
			if err != nil {
				return "", err
			}
		}
	}
	return authToken.AccessToken, nil
}
