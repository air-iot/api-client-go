package api_client_go

import (
	"context"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/apicontext"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/api-client-go/v4/warning"
	cErrors "github.com/air-iot/errors"
	"github.com/air-iot/json"
)

// QueryWarn 查询
func (c *Client) QueryWarn(ctx context.Context, projectId, token, archive string, query interface{}, result interface{}) (int, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.NewMsg("序列化查询参数为空, %s", err)
	}
	cli, err := c.WarningClient.GetWarnServiceClient()
	if err != nil {
		return 0, errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&warning.QueryWarningRequest{Query: bts, Archive: archive})
	if err != nil {
		return 0, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return 0, cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return 0, errors.NewMsg("解析请求结果错误, %s", err)
	}
	return int(res.Count), nil
}

func (c *Client) GetWarn(ctx context.Context, projectId, archive, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.NewMsg("id为空")
	}
	cli, err := c.WarningClient.GetWarnServiceClient()
	if err != nil {
		return nil, errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&warning.GetOrDeleteWarningRequest{Id: id, Archive: archive})
	if err != nil {
		return nil, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return nil, cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	if result == nil {
		return res.GetResult(), nil
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return nil, errors.NewMsg("解析请求结果错误, %s", err)
	}
	return res.GetResult(), nil
}

func (c *Client) BatchCreateWarn(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.NewMsg("插入数据为空")
	}
	cli, err := c.WarningClient.GetWarnServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.NewMsg("序列化插入数据为空")
	}
	res, err := cli.BatchCreate(
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

// QueryRule 查询
func (c *Client) QueryRule(ctx context.Context, projectId string, query interface{}, result interface{}) (int, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.NewMsg("序列化查询参数为空, %s", err)
	}
	cli, err := c.WarningClient.GetRuleServiceClient()
	if err != nil {
		return 0, errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if err != nil {
		return 0, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return 0, cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return 0, errors.NewMsg("解析请求结果错误, %s", err)
	}
	return int(res.Count), nil
}

func (c *Client) CreateWarn(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.NewMsg("插入数据为空")
	}
	cli, err := c.WarningClient.GetWarnServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.NewMsg("序列化插入数据为空")
	}
	res, err := cli.Create(
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

func (c *Client) BatchCreateRule(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.NewMsg("插入数据为空")
	}
	cli, err := c.WarningClient.GetRuleServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.NewMsg("序列化插入数据为空")
	}
	res, err := cli.BatchCreate(
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

func (c *Client) DeleteRule(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.NewMsg("id为空")
	}
	cli, err := c.WarningClient.GetRuleServiceClient()
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

func (c *Client) UpdateRule(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.NewMsg("id为空")
	}
	if updateData == nil {
		return errors.NewMsg("更新数据为空")
	}

	cli, err := c.WarningClient.GetRuleServiceClient()
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

func (c *Client) ReplaceRule(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.NewMsg("id为空")
	}
	if updateData == nil {
		return errors.NewMsg("更新数据为空")
	}
	cli, err := c.WarningClient.GetRuleServiceClient()
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
