package config

type Config struct {
	Metadata map[string]string
	Services map[string]Service
	Type     KeyType
	AK       string
	SK       string
	Timeout  uint
}

type KeyType string

const (
	Tenant  KeyType = "tenant"
	Project KeyType = "project"
)

type Service struct {
	Metadata map[string]string
}

const (
	XRequestProject             = "x-request-project"
	XRequestProjectDefault      = "default"
	XRequestHeaderAuthorization = "Authorization"
	XRequestQueryAuthorization  = "token"
)
