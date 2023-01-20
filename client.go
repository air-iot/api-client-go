package api_client_go

import (
	"context"
	"fmt"
	"github.com/air-iot/json"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/core"
	"github.com/air-iot/api-client-go/v4/dataservice"
	"github.com/air-iot/api-client-go/v4/driver"
	"github.com/air-iot/api-client-go/v4/engine"
	"github.com/air-iot/api-client-go/v4/flow"
	"github.com/air-iot/api-client-go/v4/metadata"
	"github.com/air-iot/api-client-go/v4/spm"
	"github.com/air-iot/api-client-go/v4/warning"
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

	Config config.Config
	tokens sync.Map
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

func (c *Client) GetToken(projectId string) (*AuthToken, error) {
	//c.lock.Lock()
	//defer c.lock.Unlock()
	//if c.ExpiresAt > time.Now().Unix() {
	//	return c.AccessToken, nil
	//}
	//if c.appServiceClient == nil {
	//	return "", fmt.Errorf("客户端是空")
	//}
	c.lock.Lock()
	defer c.lock.Unlock()
	cli, err := c.CoreClient.GetAppServiceClient()
	if err != nil {
		return nil, fmt.Errorf("获取客户端错误,%s", err)
	}

	res, err := cli.GetToken(
		metadata.GetGrpcContext(context.Background(), map[string]string{config.XRequestProject: projectId}),
		&api.TokenRequest{Ak: c.Config.AK, Sk: c.Config.SK})
	if err != nil {
		return nil, fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return nil, fmt.Errorf("响应不成功, %s %s", res.GetInfo(),res.GetDetail())
	}
	var authToken AuthToken
	if err := json.Unmarshal(res.GetResult(), &authToken); err != nil {
		return nil, fmt.Errorf("解析 token 请求结果错误, %s", err)
	}
	c.tokens.Store(projectId, &authToken)
	return &authToken, nil
}

func (c *Client) Token(projectId string) (token string, err error) {
	var authToken *AuthToken
	tokenValue, ok := c.tokens.Load(projectId)
	if !ok {
		authToken, err = c.GetToken(projectId)
		if err != nil {
			return "", err
		}
	} else {
		authToken = tokenValue.(*AuthToken)
		if authToken.ExpiresAt <= time.Now().Unix() {
			authToken, err = c.GetToken(projectId)
			if err != nil {
				return "", err
			}
		}
	}
	return authToken.AccessToken, nil
}
