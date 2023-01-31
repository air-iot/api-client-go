package api_client_go

import (
	"context"
	"fmt"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/json"
)

func (c *Client) QueryProject(ctx context.Context, query, result interface{}) error {
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	res, err := cli.Query(ctx, &api.QueryRequest{Query: bts})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) GetProject(ctx context.Context, id string, result interface{}) error {
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	res, err := cli.Get(ctx, &api.GetOrDeleteRequest{Id: id})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) DeleteProject(ctx context.Context, id string, result interface{}) error {
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	res, err := cli.Delete(ctx, &api.GetOrDeleteRequest{Id: id})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) UpdateProject(ctx context.Context, id string, updateData, result interface{}) error {
	if id == "" {
		return fmt.Errorf("id为空")
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("marshal 更新数据为空")
	}
	res, err := cli.Update(ctx, &api.UpdateRequest{Id: id, Data: bts})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) UpdateProjectLicense(ctx context.Context, id string, updateData, _ interface{}) error {
	if id == "" {
		return fmt.Errorf("id为空")
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("marshal 更新数据为空")
	}
	res, err := cli.UpdateLicense(ctx, &api.UpdateRequest{Id: id, Data: bts})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	//if err := json.Unmarshal(res.GetResult(), result); err != nil {
	//	return fmt.Errorf("解析请求结果错误, %s", err)
	//}
	return nil
}

func (c *Client) ReplaceProject(ctx context.Context, id string, updateData, result interface{}) error {

	if id == "" {
		return fmt.Errorf("id为空")
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("marshal 更新数据为空")
	}
	res, err := cli.Replace(ctx, &api.UpdateRequest{Id: id, Data: bts})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) CreateProject(ctx context.Context, createData, result interface{}) error {
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return fmt.Errorf("marshal 插入数据为空")
	}
	res, err := cli.Create(ctx, &api.CreateRequest{Data: bts})
	if err != nil {
		return fmt.Errorf("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}
