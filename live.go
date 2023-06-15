package api_client_go

import (
	"context"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/apicontext"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/errors"
	cErrors "github.com/air-iot/errors"
	"github.com/air-iot/json"
)

func (c *Client) RtspPull(ctx context.Context, projectId string, createData interface{}) (string, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return "", errors.NewMsg("插入数据为空")
	}
	cli, err := c.LiveClient.GetLiveServiceClient()
	if err != nil {
		return "", errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return "", errors.NewMsg("序列化插入数据为空")
	}
	res, err := cli.Create(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if err != nil {
		return "", errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return "", cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	result := res.GetResult()
	if result != nil || string(result) == "" {
		return string(result), nil
	} else {
		return "", errors.NewMsg("result中path为空, %s", err)
	}
}
