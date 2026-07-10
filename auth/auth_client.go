package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/air-iot/api-client-go/v4/api"
	context2 "github.com/air-iot/api-client-go/v4/apicontext"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/core"
	"github.com/air-iot/api-client-go/v4/spm"
	"github.com/air-iot/errors"
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
			return nil, errors.Wrap(err, "获取token客户端错误")
		}
		res, err := cli.GetToken(context.Background(), &api.TokenRequest{Ak: a.cfg.AK, Sk: a.cfg.SK})
		if err != nil {
			return nil, errors.NewResErrorMsg(err, "请求token错误")
		}
		if !res.GetStatus() {
			return nil, errors.Wrap400Response(fmt.Errorf(res.GetDetail()), int(res.GetCode()), "请求token响应错误: %s", res.GetInfo())
		}
		if err := json.Unmarshal(res.GetResult(), &authToken); err != nil {
			return nil, errors.Wrap(err, "解析 token 请求结果错误")
		}
	case config.Project:
		cli, err := a.coreClient.GetAppServiceClient()
		if err != nil {
			return nil, errors.Wrap(err, "获取token客户端错误")
		}
		res, err := cli.GetToken(context2.GetGrpcContext(context.Background(), map[string]string{config.XRequestProject: a.cfg.ProjectId}), &api.TokenRequest{Ak: a.cfg.AK, Sk: a.cfg.SK})
		if err != nil {
			return nil, errors.NewResErrorMsg(err, "请求token错误")
		}
		if !res.GetStatus() {
			return nil, errors.Wrap400Response(fmt.Errorf(res.GetDetail()), int(res.GetCode()), "请求token响应错误: %s", res.GetInfo())
		}
		if err := json.Unmarshal(res.GetResult(), &authToken); err != nil {
			return nil, errors.Wrap(err, "解析 token 请求结果错误")
		}
	default:
		return nil, errors.New("未知ak、sk类型")
	}
	a.authToken = authToken
	a.hasToken = true
	return &authToken, nil
}

func (a *Client) Token() (token string, err error) {
	// API Key 认证:直接返回配置的 ApiKey,无需获取 token
	if a.cfg.AuthType == config.AuthTypeApiKey {
		return a.cfg.ApiKey, nil
	}
	a.lock.Lock()
	defer a.lock.Unlock()
	var authToken *Token
	if !a.hasToken {
		authToken, err = a.getToken()
		if err != nil {
			return "", err
		}
	} else {
		if a.authToken.ExpiresAt <= (time.Now().Unix() + a.cfg.ExpirePrecision) {
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

// Authorization 返回 HTTP Authorization 头的完整值。
// token 认证为 "Bearer <token>";API Key 认证直接返回 ApiKey。
func (a *Client) Authorization() (string, error) {
	token, err := a.Token()
	if err != nil {
		return "", err
	}
	if a.cfg.AuthType == config.AuthTypeApiKey {
		return token, nil
	}
	return fmt.Sprintf("Bearer %s", token), nil
}
