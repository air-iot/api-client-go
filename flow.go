package api_client_go

import (
	"context"
	"fmt"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/metadata"
	"github.com/air-iot/json"
)

// CreateFlowTask FlowTask
func (c *Client) CreateFlowTask(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.FlowClient.GetFlowTaskServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token := ""
	bts, err := json.Marshal(createData)
	if err != nil {
		return fmt.Errorf("marshal 插入数据为空")
	}
	res, err := cli.Create(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.CreateRequest{Data: bts})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(),res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) GetFlowTask(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.FlowClient.GetFlowTaskServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token := ""
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(),res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

// QueryFlow Flow
func (c *Client) QueryFlow(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	cli, err := c.FlowClient.GetFlowServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token := ""
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.QueryRequest{Query: bts})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(),res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) GetFlow(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.FlowClient.GetFlowServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token := ""
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(),res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) UpdateFlow(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}
	cli, err := c.FlowClient.GetFlowServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token := ""
	bts, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("marshal 更新数据为空")
	}
	res, err := cli.Update(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.UpdateRequest{Id: id, Data: bts})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(),res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

// CreateFlowTriggerRecord FlowTriggerRecord
func (c *Client) CreateFlowTriggerRecord(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.FlowClient.GetFlowTriggerRecordServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token := ""
	bts, err := json.Marshal(createData)
	if err != nil {
		return fmt.Errorf("marshal 插入数据为空")
	}
	res, err := cli.Create(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.CreateRequest{Data: bts})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(),res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}
