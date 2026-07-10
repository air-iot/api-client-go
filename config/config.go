package config

import "time"

type ConnectType string

const (
	Local ConnectType = "local"
	Grpc  ConnectType = "grpc"
)

type Config struct {
	LiteMode    bool               `json:"liteMode"`
	Gateway     string             `json:"gateway"`
	GatewayGrpc string             `json:"gatewayGrpc"`
	EtcdConfig  string             `json:"etcdConfig"`
	Metadata    map[string]string  `json:"metadata"`
	Services    map[string]Service `json:"services"`
	Type        KeyType            `json:"type"`
	ProjectId   string             `json:"projectId"`
	AK          string             `json:"ak"`
	SK          string             `json:"sk"`
	// AuthType 认证方式,默认空(token),配置为 apiKey 时使用 API Key 认证
	AuthType AuthType `json:"authType"`
	// ApiKey API Key,当 AuthType 为 apiKey 时使用,作为请求 Authorization 的值,无需获取 token
	ApiKey          string        `json:"apiKey"`
	Timeout         uint          `json:"timeout"`
	KeepAlive       bool          `json:"keepAlive"`
	MaxIdleConns    int           `json:"maxIdleConns"`
	IdleConnTimeout time.Duration `json:"idleConnTimeout"`
	Limit           int           `json:"limit"`
	Debug           bool          `json:"debug"`
	Service         struct {
		//Enable bool          `json:"enable"`
		Expire time.Duration `json:"expire"`
	} `json:"service"`
	ExpirePrecision int64 `json:"expirePrecision"`
}

type KeyType string

const (
	Tenant  KeyType = "tenant"
	Project KeyType = "project"
)

// AuthType 认证方式
type AuthType string

const (
	// AuthTypeToken token 认证(默认,空值),通过 AK/SK 获取 token
	AuthTypeToken AuthType = ""
	// AuthTypeApiKey API Key 认证,配置 ApiKey 后无需获取 token,
	// 请求时 Authorization 的值即为 ApiKey
	AuthTypeApiKey AuthType = "apiKey"
)

type Service struct {
	Metadata map[string]string `json:"metadata"`
	GrpcAddr string            `json:"grpcAddr"`
}

const (
	XRequestProject             = "x-request-project"
	XRequestProjectDefault      = "default"
	XRequestHeaderAuthorization = "Authorization"
	XRequestQueryAuthorization  = "token"
)
