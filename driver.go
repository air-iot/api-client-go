package api_client_go

import (
	"context"
	"fmt"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/metadata"
	"github.com/air-iot/json"
)

// BatchCommand 批量执行指令
func (c *Client) BatchCommand(ctx context.Context, projectId string, data interface{}, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if data == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.DriverClient.GetDriverServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal 插入数据为空")
	}
	token := ""
	res, err := cli.BatchCommand(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.CreateRequest{Data: bts})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s", res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

// ChangeCommand 执行指令
func (c *Client) ChangeCommand(ctx context.Context, projectId, id string, data, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.DriverClient.GetDriverServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token := ""
	if data == nil {
		return fmt.Errorf("更新数据为空")
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal 更新数据为空")
	}
	res, err := cli.ChangeCommand(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.UpdateRequest{Id: id, Data: bts})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s", res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}
