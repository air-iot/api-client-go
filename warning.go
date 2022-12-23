package api_client_go

import (
	"context"
	"fmt"
	"strconv"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/metadata"
	"github.com/air-iot/api-client-go/v4/warning"
	"github.com/air-iot/json"
)

// QueryWarn 查询
func (c *Client) QueryWarn(ctx context.Context, projectId, token, archive string, query interface{}, result interface{}) (int, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, fmt.Errorf("序列化查询参数为空, %s", err)
	}
	cli, err := c.WarningClient.GetWarnServiceClient()
	if err != nil {
		return 0, fmt.Errorf("获取客户端错误,%s", err)
	}
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&warning.QueryWarningRequest{Query: bts, Archive: archive})
	if err != nil {
		return 0, fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return 0, fmt.Errorf("响应不成功, %s", res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return 0, fmt.Errorf("解析请求结果错误, %s", err)
	}
	count := 0
	if res.GetDetail() != "" {
		count, err = strconv.Atoi(res.GetDetail())
		if err != nil {
			count = 0
			//return 0,fmt.Errorf("查询结果数量(%s)转数字失败, %s",res.GetDetail(), err)
		}
	}
	return count, nil
}

func (c *Client) GetWarn(ctx context.Context, projectId, archive, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.WarningClient.GetWarnServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token := ""
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&warning.GetOrDeleteWarningRequest{Id: id, Archive: archive})
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

// QueryRule 查询
func (c *Client) QueryRule(ctx context.Context, projectId string, query interface{}, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	token := ""
	cli, err := c.WarningClient.GetRuleServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.QueryRequest{Query: bts})
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

func (c *Client) CreateWarn(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.WarningClient.GetWarnServiceClient()
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
		return fmt.Errorf("响应不成功, %s", res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}
