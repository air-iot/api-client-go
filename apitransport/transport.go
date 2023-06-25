package apitransport

import (
	"context"
	netHttp "net/http"

	"github.com/go-kratos/kratos/v2/transport/http"
)

type apiClientTransportKey struct{}

// Transport is an HTTP transport.
type Transport struct {
	ReqHeader    map[string]string
	replyHeader  netHttp.Header
	request      *http.Request
	pathTemplate string
}

func NewClientContext(ctx context.Context, tt *Transport) context.Context {
	return context.WithValue(ctx, apiClientTransportKey{}, tt)
}

// FromClientContext returns the Transport value stored in ctx, if any.
func FromClientContext(ctx context.Context) (tt *Transport, ok bool) {
	tt, ok = ctx.Value(apiClientTransportKey{}).(*Transport)
	return
}
