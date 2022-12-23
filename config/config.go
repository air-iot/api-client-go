package config

type Config struct {
	Metadata map[string]string
	Services map[string]Service
	AK       string
	SK       string
	Timeout  uint
}

type Service struct {
	Metadata map[string]string
}

const (
	XRequestProject             = "x-request-project"
	XRequestProjectDefault      = "default"
	XRequestHeaderAuthorization = "Authorization"
	XRequestQueryAuthorization  = "token"
)
