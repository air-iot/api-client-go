package api_client_go

import (
	"context"
	cErrors "github.com/air-iot/errors"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/api-client-go/v4/metadata"
	"github.com/air-iot/json"
)

// BatchCommand 批量执行指令
func (c *Client) BatchCommand(ctx context.Context, projectId string, data interface{}, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if data == nil {
		return errors.NewMsg("插入数据为空")
	}
	cli, err := c.DriverClient.GetDriverServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return errors.NewMsg("marshal 插入数据为空")
	}
	res, err := cli.BatchCommand(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return errors.NewMsg("解析请求结果错误, %s", err)
	}
	return nil
}

// ChangeCommand 执行指令
func (c *Client) ChangeCommand(ctx context.Context, projectId, id string, data, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.NewMsg("id为空")
	}
	cli, err := c.DriverClient.GetDriverServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	if data == nil {
		return errors.NewMsg("更新数据为空")
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return errors.NewMsg("marshal 更新数据为空")
	}
	res, err := cli.ChangeCommand(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return errors.NewMsg("解析请求结果错误, %s", err)
	}
	return nil
}
