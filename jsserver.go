package api_client_go

import (
	"context"
	"fmt"

	"github.com/air-iot/api-client-go/v4/jsserver"
	"github.com/air-iot/errors"
	"github.com/air-iot/json"
)

func (c *Client) RunJsScript(ctx context.Context, variables interface{}, script string) ([]byte, error) {
	cli, err := c.JsServerClient.GetScriptClient()
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(variables)
	if err != nil {
		return nil, errors.Wrap(err, "序列化请求数据错误")
	}
	res, err := cli.Run(ctx, &jsserver.Request{
		Content: script,
		Params:  b,
	})
	if err != nil {
		return nil, errors.Wrap(err, "请求错误")
	}
	if res.GetStatus() {
		return res.GetResult(), nil
	} else {
		return nil, errors.Wrap(fmt.Errorf("%s", res.GetDetail()), res.GetInfo())
	}
}
