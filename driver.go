package api_client_go

import (
	"context"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/apicontext"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/driver"
	"github.com/air-iot/api-client-go/v4/errors"
	cErrors "github.com/air-iot/errors"
	"github.com/air-iot/json"
)

type DriverWriteTag struct {
	Table     string      `json:"table"`
	TableData string      `json:"tableData"`
	ID        string      `json:"id" bson:"id"`
	Params    interface{} `json:"params" bson:"params"`
}

type DriverBatchWriteTag struct {
	TableId      string      `json:"tableId"`
	TableDataIds []string    `json:"tableDataIds"`
	ID           string      `json:"id" bson:"id"`
	Query        string      `json:"query"`
	Type         string      `json:"type"`
	Params       interface{} `json:"params" bson:"params"`
}

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
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
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
	//if id == "" {
	//	return errors.NewMsg("id为空")
	//}
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
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
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

// DriverWriteTag 执行写数据点
func (c *Client) DriverWriteTag(ctx context.Context, projectId string, data *DriverWriteTag, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
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
	res, err := cli.WriteTag(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
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

// DriverBatchWriteTag 执行写数据点
func (c *Client) DriverBatchWriteTag(ctx context.Context, projectId string, data *DriverBatchWriteTag, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
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
	res, err := cli.BatchWriteTag(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
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

// HttpProxy 驱动代理接口
func (c *Client) HttpProxy(ctx context.Context, projectId, typeId, groupId string, headers, data, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.DriverClient.GetDriverServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	if data == nil {
		return errors.NewMsg("数据为空")
	}

	var headersBytes []byte
	if headers != nil {
		headersBytes, err = json.Marshal(headers)
		if err != nil {
			return errors.NewMsg("headers序列化失败")
		}
	}

	bts, err := json.Marshal(data)
	if err != nil {
		return errors.NewMsg("数据为空")
	}
	res, err := cli.HttpProxy(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&driver.ClientHttpProxyRequest{
			ProjectId: projectId,
			Type:      typeId,
			GroupId:   groupId,
			Headers:   headersBytes,
			Data:      bts,
		})
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

func (c *Client) QueryDriverInstance(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.NewMsg("序列化查询参数为空, %s", err)
	}
	cli, err := c.DriverClient.GetDriverInstanceServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return errors.NewErrorMsg(errors.NewMsg("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return errors.NewMsg("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) GetDriverInstance(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.NewMsg("id为空")
	}
	cli, err := c.DriverClient.GetDriverInstanceServiceClient()
	if err != nil {
		return nil, errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) DeleteDriverInstance(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.NewMsg("id为空")
	}
	cli, err := c.DriverClient.GetDriverInstanceServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}

	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return errors.NewErrorMsg(errors.NewMsg("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return errors.NewMsg("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) UpdateDriverInstance(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.NewMsg("id为空")
	}
	if updateData == nil {
		return errors.NewMsg("更新数据为空")
	}

	cli, err := c.DriverClient.GetDriverInstanceServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}

	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.NewMsg("marshal 更新数据为空")
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return errors.NewErrorMsg(errors.NewMsg("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return errors.NewMsg("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) ReplaceDriverInstance(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.NewMsg("id为空")
	}
	if updateData == nil {
		return errors.NewMsg("更新数据为空")
	}
	cli, err := c.DriverClient.GetDriverInstanceServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.NewMsg("marshal 更新数据为空")
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return errors.NewErrorMsg(errors.NewMsg("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return errors.NewMsg("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) CreateDriverInstance(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.NewMsg("插入数据为空")
	}
	cli, err := c.DriverClient.GetDriverInstanceServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.NewMsg("marshal 插入数据为空")
	}
	res, err := cli.Create(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return errors.NewErrorMsg(errors.NewMsg("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return errors.NewMsg("解析请求结果错误, %s", err)
	}
	return nil
}
