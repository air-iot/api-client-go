package api_client_go

import (
	"context"
	"fmt"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/metadata"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/json"
)

func (c *Client) QueryProject(ctx context.Context, query, result interface{}) error {
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误, %s", err)
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.NewMsg("序列化查询参数为空, %s", err)
	}
	token, err := c.Token("")
	if err != nil {
		return errors.NewMsg("查询token错误, %s", err)
	}
	res, err := cli.Query(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestHeaderAuthorization: token}), &api.QueryRequest{Query: bts})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return errors.NewErrorMsg(fmt.Errorf("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return errors.NewMsg("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) GetProject(ctx context.Context, id string, result interface{}) error {
	if id == "" {
		return errors.NewMsg("id为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	token, err := c.Token("")
	if err != nil {
		return errors.NewMsg("查询token错误, %s", err)
	}
	res, err := cli.Get(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestHeaderAuthorization: token}), &api.GetOrDeleteRequest{Id: id})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return errors.NewErrorMsg(fmt.Errorf("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return errors.NewMsg("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) DeleteProject(ctx context.Context, id string, result interface{}) error {
	if id == "" {
		return errors.NewMsg("id为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	token, err := c.Token("")
	if err != nil {
		return errors.NewMsg("查询token错误, %s", err)
	}
	res, err := cli.Delete(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestHeaderAuthorization: token}), &api.GetOrDeleteRequest{Id: id})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return errors.NewErrorMsg(fmt.Errorf("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return errors.NewMsg("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) UpdateProject(ctx context.Context, id string, updateData, result interface{}) error {
	if id == "" {
		return errors.NewMsg("id为空")
	}
	if updateData == nil {
		return errors.NewMsg("更新数据为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.NewMsg("序列化更新数据为空")
	}
	token, err := c.Token("")
	if err != nil {
		return errors.NewMsg("查询token错误, %s", err)
	}
	res, err := cli.Update(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestHeaderAuthorization: token}), &api.UpdateRequest{Id: id, Data: bts})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return errors.NewErrorMsg(fmt.Errorf("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return errors.NewMsg("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) UpdateProjectLicense(ctx context.Context, id string, updateData, _ interface{}) error {
	if id == "" {
		return errors.NewMsg("id为空")
	}
	if updateData == nil {
		return errors.NewMsg("更新数据为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.NewMsg("序列化更新数据为空")
	}
	token, err := c.Token("")
	if err != nil {
		return errors.NewMsg("查询token错误, %s", err)
	}
	res, err := cli.UpdateLicense(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestHeaderAuthorization: token}), &api.UpdateRequest{Id: id, Data: bts})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return errors.NewErrorMsg(fmt.Errorf("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	//if err := json.Unmarshal(res.GetResult(), result); err != nil {
	//	return fmt.Errorf("解析请求结果错误, %s", err)
	//}
	return nil
}

func (c *Client) ReplaceProject(ctx context.Context, id string, updateData, result interface{}) error {

	if id == "" {
		return errors.NewMsg("id为空")
	}
	if updateData == nil {
		return errors.NewMsg("更新数据为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.NewMsg("序列化更新数据为空")
	}
	token, err := c.Token("")
	if err != nil {
		return errors.NewMsg("查询token错误, %s", err)
	}
	res, err := cli.Replace(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestHeaderAuthorization: token}), &api.UpdateRequest{Id: id, Data: bts})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return errors.NewErrorMsg(fmt.Errorf("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return errors.NewMsg("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) CreateProject(ctx context.Context, createData, result interface{}) error {
	if createData == nil {
		return errors.NewMsg("插入数据为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.NewMsg("序列化插入数据为空")
	}
	token, err := c.Token("")
	if err != nil {
		return errors.NewMsg("查询token错误, %s", err)
	}
	res, err := cli.Create(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestHeaderAuthorization: token}), &api.CreateRequest{Data: bts})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return errors.NewErrorMsg(fmt.Errorf("响应不成功, %s", res.GetDetail()), res.GetInfo())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return errors.NewMsg("解析请求结果错误, %s", err)
	}
	return nil
}
