package api_client_go

import (
	"context"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/dataservice"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/api-client-go/v4/metadata"
	"github.com/air-iot/json"
)

type ProxyResult struct {
	Code    int32  `json:"code"`
	Headers []byte `json:"headers"`
	Body    []byte `json:"body"`
}

func (c *Client) QueryDataGroup(ctx context.Context, projectId string, query, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.NewMsg("序列化查询参数为空, %s", err)
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return 0, errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if err != nil {
		return 0, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return 0, errors.NewErrorMsg(errors.NewMsg("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return 0, errors.NewMsg("解析请求结果错误, %s", err)
	}
	return res.GetCount(), nil
}

func (c *Client) CreateDataGroups(ctx context.Context, projectId string, createData, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return 0, errors.NewMsg("插入数据为空")
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return 0, errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return 0, errors.NewMsg("marshal 插入数据为空")
	}
	res, err := cli.CreateMany(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{
			Data: bts,
		})
	if err != nil {
		return 0, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return 0, errors.NewErrorMsg(errors.NewMsg("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	//if err := json.Unmarshal(res.GetResult(), result); err != nil {
	//	return 0, errors.NewMsg("解析请求结果错误, %s", err)
	//}
	return res.GetCount(), nil
}

func (c *Client) QueryDataInterface(ctx context.Context, projectId string, query, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.NewMsg("序列化查询参数为空, %s", err)
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
	if err != nil {
		return 0, errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if err != nil {
		return 0, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return 0, errors.NewErrorMsg(errors.NewMsg("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return 0, errors.NewMsg("解析请求结果错误, %s", err)
	}
	return res.GetCount(), nil
}

func (c *Client) CreateDataInterfaces(ctx context.Context, projectId string, createData, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return 0, errors.NewMsg("插入数据为空")
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
	if err != nil {
		return 0, errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return 0, errors.NewMsg("marshal 插入数据为空")
	}
	res, err := cli.CreateMany(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{
			Data: bts,
		})
	if err != nil {
		return 0, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return 0, errors.NewErrorMsg(errors.NewMsg("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	//if err := json.Unmarshal(res.GetResult(), result); err != nil {
	//	return 0, errors.NewMsg("解析请求结果错误, %s", err)
	//}
	return res.GetCount(), nil
}

func (c *Client) DataInterfaceProxy(ctx context.Context, projectId, key string, data map[string]interface{}) (*ProxyResult, error) {
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
	res, err := cli.Proxy(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&dataservice.Request{Key: key, Data: bts})
	if err != nil {
		return nil, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return nil, errors.NewErrorMsg(errors.NewMsg("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	return &ProxyResult{
		Code:    res.GetHttpCode(),
		Headers: res.GetHeaders(),
		Body:    res.GetResult(),
	}, nil
	//if res.GetResult() == nil || len(res.GetResult()) == 0 {
	//	return res.GetResult(), nil
	//}
	//if err := json.Unmarshal(res.GetResult(), result); err != nil {
	//	return nil, errors.NewMsg("解析请求结果错误, %s", err)
	//}
	//return res.GetResult(), nil
}
