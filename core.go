package api_client_go

import (
	"context"
	"fmt"

	"github.com/air-iot/json"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/core"
	"github.com/air-iot/api-client-go/v4/metadata"
)

func (c *Client) GetFileLicense(ctx context.Context, result interface{}) error {
	cli, err := c.CoreClient.GetLicenseServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	res, err := cli.GetFileLicense(ctx, &api.QueryRequest{})
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

func (c *Client) UseLicense(ctx context.Context, projectId string, result interface{}) error {
	cli, err := c.CoreClient.GetLicenseServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	res, err := cli.UseLicense(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &api.QueryRequest{})
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

func (c *Client) GetDriverLicense(ctx context.Context, projectId, driverId string, result interface{}) error {
	cli, err := c.CoreClient.GetLicenseServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	res, err := cli.GetDriverLicense(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &api.GetOrDeleteRequest{
		Id: driverId,
	})
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

func (c *Client) GetCurrentUserInfo(ctx context.Context, projectId, token string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if token == "" {
		return fmt.Errorf("token is empty")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	res, err := cli.GetCurrentUserInfo(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&core.LoginUserRequest{})
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

func (c *Client) QueryUser(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}

	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.QueryRequest{Query: bts})
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

func (c *Client) GetUser(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) DeleteUser(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Delete(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) UpdateUser(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}

	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
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
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) ReplaceUser(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("marshal 更新数据为空")
	}
	res, err := cli.Replace(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.UpdateRequest{Id: id, Data: bts})
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

func (c *Client) CreateLog(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}

	cli, err := c.CoreClient.GetLogServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
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
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) CreateUser(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
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
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) QueryTableSchema(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.QueryRequest{Query: bts})
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

func (c *Client) QueryTableSchemaDeviceByDriverAndGroup(ctx context.Context, projectId, driverId, groupId string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if groupId == "" {
		return fmt.Errorf("groupId is empty")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	//token, err := c.Token(projectId)
	//if err != nil {
	//	return fmt.Errorf("查询token错误, %s", err)
	//}
	token := ""
	res, err := cli.QueryDeviceByDriverAndGroup(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&core.GetDeviceRequest{Driver: driverId, Group: groupId})
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

func (c *Client) QueryEmulator(ctx context.Context, projectId string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	//token, err := c.Token(projectId)
	//if err != nil {
	//	return fmt.Errorf("查询token错误, %s", err)
	//}
	token := ""
	res, err := cli.QueryEmulator(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.QueryRequest{})
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

func (c *Client) GetTableSchema(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) DeleteTableSchema(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Delete(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) UpdateTableSchema(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
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
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) ReplaceTableSchema(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("marshal 更新数据为空")
	}
	res, err := cli.Replace(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.UpdateRequest{Id: id, Data: bts})
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

func (c *Client) CreateTableSchema(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
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
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) QueryTableRecord(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	cli, err := c.CoreClient.GetTableRecordServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.QueryRequest{Query: bts})
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

func (c *Client) GetTableRecord(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetTableRecordServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) DeleteTableRecord(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetTableRecordServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Delete(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) UpdateTableRecord(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}
	cli, err := c.CoreClient.GetTableRecordServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
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
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) ReplaceTableRecord(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}
	cli, err := c.CoreClient.GetTableRecordServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("marshal 更新数据为空")
	}
	res, err := cli.Replace(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.UpdateRequest{Id: id, Data: bts})
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

func (c *Client) CreateTableRecord(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.CoreClient.GetTableRecordServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
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
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) QueryTableData(ctx context.Context, projectId, tableName string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return fmt.Errorf("表为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&core.QueryDataRequest{
			Table: tableName,
			Query: bts,
		})
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

func (c *Client) QueryTableDataByTableId(ctx context.Context, projectId, tableId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableId == "" {
		return fmt.Errorf("记录id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	res, err := cli.QueryByTableId(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&core.QueryDataRequest{
			Table: tableId,
			Query: bts,
		})
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

func (c *Client) GetTableData(ctx context.Context, projectId, tableName, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return fmt.Errorf("表为空")
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&core.GetOrDeleteDataRequest{Table: tableName, Id: id})
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

func (c *Client) DeleteTableData(ctx context.Context, projectId, tableName, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return fmt.Errorf("表为空")
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Delete(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&core.GetOrDeleteDataRequest{Table: tableName, Id: id})
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

func (c *Client) DeleteManyTableData(ctx context.Context, projectId, tableName string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	if tableName == "" {
		return fmt.Errorf("表为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.DeleteMany(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&core.QueryDataRequest{Table: tableName, Query: bts})
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

func (c *Client) UpdateTableData(ctx context.Context, projectId, tableName, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return fmt.Errorf("表为空")
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("marshal 更新数据为空")
	}
	res, err := cli.Update(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&core.UpdateDataRequest{Table: tableName, Id: id, Data: bts})
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

func (c *Client) ReplaceTableData(ctx context.Context, projectId, tableName, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}
	if tableName == "" {
		return fmt.Errorf("表为空")
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("marshal 更新数据为空")
	}
	res, err := cli.Replace(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&core.UpdateDataRequest{Table: tableName, Id: id, Data: bts})
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

func (c *Client) CreateTableData(ctx context.Context, projectId, tableName string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return fmt.Errorf("表为空")
	}
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return fmt.Errorf("marshal 插入数据为空")
	}
	res, err := cli.Create(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&core.CreateDataRequest{
			Table: tableName,
			Data:  bts,
		})
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

func (c *Client) CreateManyTableData(ctx context.Context, projectId, tableName string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return fmt.Errorf("表为空")
	}
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return fmt.Errorf("marshal 插入数据为空")
	}
	res, err := cli.CreateMany(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&core.CreateDataRequest{
			Table: tableName,
			Data:  bts,
		})
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

func (c *Client) CreateMessage(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.CoreClient.GetMessageServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
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
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) GetLog(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetLogServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) PostLatest(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.CoreClient.GetDataQueryServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return fmt.Errorf("marshal 插入数据为空")
	}
	res, err := cli.PostLatest(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.CreateRequest{Data: bts})
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

func (c *Client) GetQuery(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	cli, err := c.CoreClient.GetDataQueryServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.GetQuery(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.QueryRequest{Query: bts})
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

func (c *Client) PostQuery(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.CoreClient.GetDataQueryServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return fmt.Errorf("marshal 插入数据为空")
	}
	res, err := cli.PostQuery(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.CreateRequest{Data: bts})
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

func (c *Client) QueryRole(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	cli, err := c.CoreClient.GetRoleServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.QueryRequest{Query: bts})
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

func (c *Client) AdminRoleCheck(ctx context.Context, projectId, token string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if token == "" {
		return fmt.Errorf("无Token认证信息")
	}
	cli, err := c.CoreClient.GetRoleServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	res, err := cli.AdminRoleCheck(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.EmptyRequest{})
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

func (c *Client) GetRole(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetRoleServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) FindTableDataDeptByDeptIDs(ctx context.Context, projectId string, ids map[string]interface{}, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	bts, err := json.Marshal(ids)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	res, err := cli.FindTableDataDeptByDeptIDs(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.CreateRequest{
			Data: bts,
		})
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

func (c *Client) UpdateManyTableData(ctx context.Context, projectId, tableName string, query, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	if tableName == "" {
		return fmt.Errorf("表为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	btsUpdate, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("marshal 更新数据为空")
	}
	res, err := cli.UpdateMany(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&core.MultiUpdateDataRequest{Table: tableName, Query: bts, Data: btsUpdate})
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

func (c *Client) GetWarningFilterIDs(ctx context.Context, projectId, token string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	if token == "" {
		return fmt.Errorf("无Token认证信息")
	}
	res, err := cli.GetWarningFilterIDs(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.EmptyRequest{})
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

func (c *Client) QueryCatalog(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	cli, err := c.CoreClient.GetCatalogServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.QueryRequest{Query: bts})
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

func (c *Client) GetCatalog(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetCatalogServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) QueryDept(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	cli, err := c.CoreClient.GetDeptServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.QueryRequest{Query: bts})
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

func (c *Client) GetDept(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetDeptServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) QuerySetting(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	cli, err := c.CoreClient.GetSettingServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.QueryRequest{Query: bts})
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

func (c *Client) QueryApp(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	cli, err := c.CoreClient.GetAppServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.QueryRequest{Query: bts})
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

func (c *Client) GetApp(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetAppServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) QuerySystemVariable(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("序列化查询参数为空, %s", err)
	}
	cli, err := c.CoreClient.GetSystemVariableServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}

	res, err := cli.Query(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.QueryRequest{Query: bts})
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

func (c *Client) GetSystemVariable(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetSystemVariableServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Get(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) DeleteSystemVariable(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	cli, err := c.CoreClient.GetSystemVariableServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	res, err := cli.Delete(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.GetOrDeleteRequest{Id: id})
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

func (c *Client) UpdateSystemVariable(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}

	cli, err := c.CoreClient.GetSystemVariableServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
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
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}

func (c *Client) ReplaceSystemVariable(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return fmt.Errorf("id为空")
	}
	if updateData == nil {
		return fmt.Errorf("更新数据为空")
	}
	cli, err := c.CoreClient.GetSystemVariableServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("marshal 更新数据为空")
	}
	res, err := cli.Replace(
		metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.UpdateRequest{Id: id, Data: bts})
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

func (c *Client) CreateSystemVariable(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return fmt.Errorf("插入数据为空")
	}
	cli, err := c.CoreClient.GetSystemVariableServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return fmt.Errorf("查询token错误, %s", err)
	}
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
		return fmt.Errorf("响应不成功, %s %s", res.GetInfo(), res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return fmt.Errorf("解析请求结果错误, %s", err)
	}
	return nil
}
