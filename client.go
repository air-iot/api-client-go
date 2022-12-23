package api_client_go

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/air-iot/json"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"

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
	tokens map[string]*AuthToken
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
			tokens:            map[string]*AuthToken{},
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
	lock sync.RWMutex

	projectId        string
	config           config.Config
	appServiceClient core.AppServiceClient

	TokenType   string `json:"tokenType"`
	ExpiresAt   int64  `json:"expiresAt"`
	AccessToken string `json:"accessToken"`
}

func (a *AuthToken) GetToken() (string, error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.ExpiresAt > time.Now().Unix() {
		return a.AccessToken, nil
	}
	if a.appServiceClient == nil {
		return "", fmt.Errorf("客户端是空")
	}
	res, err := a.appServiceClient.GetToken(metadata.GetGrpcContext(context.Background(), map[string]string{config.XRequestProject: a.projectId}), &api.TokenRequest{Ak: a.config.AK, Sk: a.config.SK})
	if err != nil {
		return "", fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return "", fmt.Errorf("响应不成功, %s", res.GetInfo())
	}
	var authToken AuthToken
	if err := json.Unmarshal(res.GetResult(), &authToken); err != nil {
		return "", fmt.Errorf("解析 token 请求结果错误, %s", err)
	}
	a.TokenType = authToken.TokenType
	a.ExpiresAt = authToken.ExpiresAt
	a.AccessToken = authToken.AccessToken
	return authToken.AccessToken, nil
}

func (c *Client) Token(projectId string) (string, error) {
	c.lock.Lock()
	srv, ok := c.tokens[projectId]
	if !ok {
		cli, err := c.CoreClient.GetAppServiceClient()
		if err != nil {
			return "", fmt.Errorf("token客户端是空,%s", err)
		}
		srv = &AuthToken{
			projectId:        projectId,
			config:           c.Config,
			appServiceClient: cli,
		}
		c.tokens[projectId] = srv
	}
	c.lock.Unlock()
	token, err := srv.GetToken()
	if err != nil {
		return "", err
	}
	return token, nil
}
