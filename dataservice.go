package api_client_go

import (
	"context"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/dataservice"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/api-client-go/v4/metadata"
	"github.com/air-iot/json"
)

func (c *Client) DataInterfaceProxy(ctx context.Context, projectId, key string, data map[string]interface{}, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if key == "" {
		return nil, errors.NewMsg("key为空")
	}
	if data == nil {
		return nil, errors.NewMsg("请求数据为空")
	}
	cli, err := c.DataServiceClient.GetDataServiceClient()
	if err != nil {
		return nil, errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return nil, errors.NewMsg("序列化请求数据错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return nil, errors.NewMsg("查询token错误, %s", err)
	}
	res, err := cli.Proxy(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&dataservice.Request{Key: key, Data: bts})
	if err != nil {
		return nil, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return nil, errors.NewErrorMsg(errors.NewMsg("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if result == nil {
		return res.GetResult(), nil
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return nil, errors.NewMsg("解析请求结果错误, %s", err)
	}
	return res.GetResult(), nil
}
