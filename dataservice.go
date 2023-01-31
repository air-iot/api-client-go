package api_client_go

import (
	"context"
	"fmt"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/dataservice"
	"github.com/air-iot/api-client-go/v4/metadata"
	"github.com/air-iot/json"
)

func (c *Client) DataInterfaceProxy(ctx context.Context, projectId, key string, data map[string]interface{}, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if key == "" {
		return nil, fmt.Errorf("key is empty")
	}
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}
	cli, err := c.DataServiceClient.GetDataServiceClient()
	if err != nil {
		return nil, fmt.Errorf("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal data is nil")
	}
	token := ""
	res, err := cli.Proxy(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&dataservice.Request{Key: key, Data: bts})
	if err != nil {
		return nil, fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return nil, fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if result == nil {
		return res.GetResult(), nil
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return nil, fmt.Errorf("解析请求结果错误, %s", err)
	}
	return res.GetResult(), nil
}
