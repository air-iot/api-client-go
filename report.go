package api_client_go

import (
	"context"
	"strconv"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/api-client-go/v4/metadata"
	"github.com/air-iot/json"
)

// QueryReport 查询
func (c *Client) QueryReport(ctx context.Context, projectId string, query interface{}, result interface{}) (int, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.NewMsg("序列化查询参数为空, %s", err)
	}
	cli, err := c.ReportClient.GetReportServiceClient()
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
	count := 0
	if res.GetDetail() != "" {
		count, err = strconv.Atoi(res.GetDetail())
		if err != nil {
			count = 0
			//return 0,errors.NewMsg( "查询结果数量(%s)转数字失败, %s",res.GetDetail(), err)
		}
	}
	return count, nil
}

func (c *Client) GetReport(ctx context.Context, projectId,  id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.NewMsg("id为空")
	}
	cli, err := c.ReportClient.GetReportServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
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

func (c *Client) BatchCreateReport(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.NewMsg("插入数据为空")
	}
	cli, err := c.ReportClient.GetReportServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.NewMsg("序列化插入数据为空")
	}
	res, err := cli.BatchCreate(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
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

// QueryReportCopy 查询
func (c *Client) QueryReportCopy(ctx context.Context, projectId string, query interface{}, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.NewMsg("序列化查询参数为空, %s", err)
	}
	cli, err := c.ReportClient.GetReportCopyServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
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

func (c *Client) CreateReport(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.NewMsg("插入数据为空")
	}
	cli, err := c.ReportClient.GetReportServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.NewMsg("序列化插入数据为空")
	}
	res, err := cli.Create(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
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

func (c *Client) BatchCreateReportCopy(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.NewMsg("插入数据为空")
	}
	cli, err := c.ReportClient.GetReportCopyServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.NewMsg("序列化插入数据为空")
	}
	res, err := cli.BatchCreate(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
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
