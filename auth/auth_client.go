package auth

import (
	"context"
	"sync"
	"time"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/core"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/api-client-go/v4/spm"
	"github.com/air-iot/json"
)

type Token struct {
	TokenType   string `json:"tokenType"`
	ExpiresAt   int64  `json:"expiresAt"`
	AccessToken string `json:"accessToken"`
}

type Client struct {
	lock sync.RWMutex

	hasToken  bool
	authToken Token

	cfg        config.Config
	spmClient  *spm.Client
	coreClient *core.Client
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (a *Client) SetClient(spmClient *spm.Client, coreClient *core.Client) {
	a.spmClient = spmClient
	a.coreClient = coreClient
}

func (a *Client) getToken() (*Token, error) {
	var authToken Token
	switch a.cfg.Type {
	case config.Tenant:
		cli, err := a.spmClient.GetUserServiceClient()
		if err != nil {
			return nil, errors.NewMsg("获取客户端错误,%s", err)
		}
		res, err := cli.GetToken(context.Background(), &api.TokenRequest{Ak: a.cfg.AK, Sk: a.cfg.SK})
		if err != nil {
			return nil, errors.NewMsg("请求错误, %s", err)
		}
		if !res.GetStatus() {
			return nil, errors.NewMsg("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
		}
		if err := json.Unmarshal(res.GetResult(), &authToken); err != nil {
			return nil, errors.NewMsg("解析 token 请求结果错误, %s", err)
		}
	case config.Project:
		cli, err := a.coreClient.GetAppServiceClient()
		if err != nil {
			return nil, errors.NewMsg("获取客户端错误,%s", err)
		}
		res, err := cli.GetToken(context.Background(), &api.TokenRequest{Ak: a.cfg.AK, Sk: a.cfg.SK})
		if err != nil {
			return nil, errors.NewMsg("请求错误, %s", err)
		}
		if !res.GetStatus() {
			return nil, errors.NewMsg("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
		}
		if err := json.Unmarshal(res.GetResult(), &authToken); err != nil {
			return nil, errors.NewMsg("解析 token 请求结果错误, %s", err)
		}
	default:
		return nil, errors.NewMsg("未知ak、sk类型")
	}
	a.authToken = authToken
	a.hasToken = true
	return &authToken, nil
}

func (a *Client) Token() (token string, err error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	var authToken *Token
	if !a.hasToken {
		authToken, err = a.getToken()
		if err != nil {
			return "", err
		}
	} else {
		if a.authToken.ExpiresAt <= time.Now().Unix() {
			authToken, err = a.getToken()
			if err != nil {
				return "", err
			}
		} else {
			authToken = &a.authToken
		}
	}
	return authToken.AccessToken, nil
}
