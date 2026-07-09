package api_client_go

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"maps"
	netHttp "net/http"
	"net/textproto"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/apicontext"
	"github.com/air-iot/api-client-go/v4/apitransport"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/core"
	"github.com/air-iot/errors"
	"github.com/air-iot/json"
	"github.com/air-iot/logger"
)

// MediaFile 媒体库文件
type MediaFile struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (c *Client) GetFileLicense(ctx context.Context, result interface{}) error {
	cli, err := c.CoreClient.GetLicenseServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.GetFileLicense(ctx, &api.QueryRequest{})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UseLicense(ctx context.Context, projectId string, result interface{}) error {
	cli, err := c.CoreClient.GetLicenseServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.UseLicense(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &api.QueryRequest{})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UploadLicense(ctx context.Context, projectId, filename string, size int, r io.Reader) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}

	cli, err := c.CoreClient.GetLicenseServiceClient()
	if err != nil {
		return err
	}
	stream, err := cli.UploadLicense(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, "filename": filename}))
	if err != nil {
		return errors.NewResErrorMsg(err, "请求错误")
	}

	//defer stream.CloseAndRecv()

	buffer := make([]byte, 1024)

	bytesReadAll := 0

	for {
		bytesRead, err := r.Read(buffer)
		if err != nil {
			return errors.Wrap(err, "读取授权文件错误")
		}
		err = stream.Send(&core.UploadFileRequest{Data: buffer[:bytesRead]})
		if err != nil {
			return errors.Wrap(err, "grpc发送授权文件错误")
		}
		bytesReadAll += bytesRead

		if bytesReadAll == size {
			//err := stream.CloseSend()
			//if err != nil {
			//	return fmt.Errorf("CloseSend错误:%s", err.Error())
			//}

			err := stream.Send(&core.UploadFileRequest{Data: []byte("down")})
			if err != nil {
				return errors.Wrap(err, "grpc发送结束标志错误")
			}

			m := new(api.Response)
			err = stream.RecvMsg(m)
			if err != nil {
				return errors.Wrap(err, "读取server响应错误")
			}
			fmt.Printf("上传文件成功，服务器响应结果:%+v\n", m)

			return nil
		}

	}
}

func (c *Client) GetDriverLicense(ctx context.Context, projectId, driverId string, result interface{}) error {
	cli, err := c.CoreClient.GetLicenseServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.GetDriverLicense(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &api.GetOrDeleteRequest{
		Id: driverId,
	})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) FindMachineCode(ctx context.Context, result interface{}) error {
	cli, err := c.CoreClient.GetLicenseServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.FindMachineCode(ctx, &api.QueryRequest{})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetCurrentUserInfo(ctx context.Context, projectId, token string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if token == "" {
		return errors.New("token is empty")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.GetCurrentUserInfo(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&core.LoginUserRequest{})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UserPermissionUpdate(ctx context.Context, projectId, token string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if token == "" {
		return errors.New("token is empty")
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.UserPermissionUpdate(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryUser(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryUserBackup(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.QueryBackup(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteManyUserBackup(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.DeleteManyBackup(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateManyUserBackup(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.CreateManyBackup(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{
			Data: bts,
		})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetUser(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) DeleteUser(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return err
	}

	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryAPIKey(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetAPIKeyServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetAPIKey(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetAPIKeyServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) GetAPIKeyByKeyID(ctx context.Context, projectId, keyID string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if keyID == "" {
		return nil, errors.New("keyId为空")
	}
	cli, err := c.CoreClient.GetAPIKeyServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.GetByKeyID(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.GetRequestName{Name: keyID})
	return parseRes(err, res, result)
}

func (c *Client) UpdateUser(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}

	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return err
	}

	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) ReplaceUser(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateLog(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}

	cli, err := c.CoreClient.GetLogServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.Create(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateUser(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.Create(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) StatsQuery(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.StatsQuery(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryTableSchema(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) RestQueryTableSchema(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	u := url.URL{Path: "/core/t/schema"}
	if query != nil {
		bts, err := json.Marshal(query)
		if err != nil {
			return errors.Wrap(err, "序列化查询参数为空")
		}
		params := url.Values{}
		params.Set("query", string(bts))
		u.RawQuery = params.Encode()
	}
	cli, err := c.CoreClient.GetRestClient()
	if err != nil {
		return err
	}
	if err := cli.Invoke(apitransport.NewClientContext(ctx, &apitransport.Transport{ReqHeader: map[string]string{config.XRequestProject: projectId}}), netHttp.MethodGet, u.RequestURI(), map[string]interface{}{}, result); err != nil {
		return errors.NewResError(err)
	}
	return nil
}

func (c *Client) QueryTableSchemaDeviceByDriverAndGroup(ctx context.Context, projectId, driverId, groupId string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if groupId == "" {
		return errors.New("实例组ID为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.QueryDeviceByDriverAndGroup(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.GetDeviceRequest{Driver: driverId, Group: groupId})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryOnlyTableSchemaDeviceByDriverAndGroup(ctx context.Context, projectId, driverId, groupId string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if groupId == "" {
		return errors.New("实例组ID为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.QueryTableDeviceByDriverAndGroup(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.GetDeviceRequest{Driver: driverId, Group: groupId})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) FindDevice(ctx context.Context, projectId, driverId, groupId, tableId, deviceId string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if groupId == "" {
		return errors.New("实例组ID为空")
	}
	if deviceId == "" {
		return errors.New("设备ID为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.FindDevice(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.GetDataDeviceRequest{Driver: driverId, Group: groupId, Table: tableId, Id: deviceId})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryEmulator(ctx context.Context, projectId string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.QueryEmulator(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetTableSchema(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) DeleteTableSchema(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateTableSchema(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) ReplaceTableSchema(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateTableSchema(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.Create(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryTableRecord(ctx context.Context, projectId string, query, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetTableRecordServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return 0, err
	}
	return res.GetCount(), nil
}

func (c *Client) GetTableRecord(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTableRecordServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) DeleteTableRecord(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTableRecordServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateTableRecord(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.CoreClient.GetTableRecordServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) ReplaceTableRecord(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.CoreClient.GetTableRecordServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateTableRecord(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetTableRecordServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.Create(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryTableData(ctx context.Context, projectId, tableName string, query, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return 0, errors.New("表为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return 0, err
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数为空")
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.QueryDataRequest{
			Table: tableName,
			Query: bts,
		})
	if _, err := parseRes(err, res, result); err != nil {
		return 0, err
	}
	return res.GetCount(), nil
}

func (c *Client) QueryTableDataByTableId(ctx context.Context, projectId, tableId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableId == "" {
		return errors.New("记录id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	res, err := cli.QueryByTableId(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.QueryDataRequest{
			Table: tableId,
			Query: bts,
		})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetTableData(ctx context.Context, projectId, tableName, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return nil, errors.New("表为空")
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.GetOrDeleteDataRequest{Table: tableName, Id: id})
	return parseRes(err, res, result)
}

func (c *Client) DeleteTableData(ctx context.Context, projectId, tableName, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return errors.New("表为空")
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.GetOrDeleteDataRequest{Table: tableName, Id: id})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteManyTableData(ctx context.Context, projectId, tableName string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	if tableName == "" {
		return errors.New("表为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.DeleteMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.QueryDataRequest{Table: tableName, Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateTableData(ctx context.Context, projectId, tableName, id string, closeRequire bool, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return errors.New("表为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.UpdateDataRequest{Table: tableName, Id: id, Data: bts, CloseRequire: closeRequire})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) ReplaceTableData(ctx context.Context, projectId, tableName, id string, closeRequire bool, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	if tableName == "" {
		return errors.New("表为空")
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.UpdateDataRequest{Table: tableName, Id: id, Data: bts, CloseRequire: closeRequire})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateTableData(ctx context.Context, projectId, tableName string, closeRequire bool, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return errors.New("表为空")
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.Create(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.CreateDataRequest{
			Table:        tableName,
			Data:         bts,
			CloseRequire: closeRequire,
		})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateManyTableData(ctx context.Context, projectId, tableName string, closeRequire bool, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return errors.New("表为空")
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.CreateMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.CreateDataRequest{
			Table:        tableName,
			Data:         bts,
			CloseRequire: closeRequire,
		})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateMessage(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetMessageServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.Create(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryMessage(ctx context.Context, projectId string, query, result interface{}) (int, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetMessageServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return 0, err
	}
	return int(res.GetCount()), nil
}

func (c *Client) GetLog(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetLogServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) QueryLog(ctx context.Context, projectId string, query, result interface{}) (int, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetLogServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return 0, err
	}
	return int(res.GetCount()), nil
}

func (c *Client) QueryApiPermissionLog(ctx context.Context, projectId string, query, result interface{}) (int, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetLogServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.QueryApiPermission(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return 0, err
	}
	return int(res.GetCount()), nil
}

func (c *Client) CreateApiPermission(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}

	cli, err := c.CoreClient.GetLogServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.CreateApiPermission(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) PostLatest(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetDataQueryServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.PostLatest(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetQuery(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetDataQueryServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.GetQuery(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) PostQuery(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetDataQueryServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.PostQuery(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryRole(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetRoleServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) AdminRoleCheck(ctx context.Context, projectId, token string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if token == "" {
		return errors.New("无Token认证信息")
	}
	cli, err := c.CoreClient.GetRoleServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.AdminRoleCheck(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.EmptyRequest{})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetRole(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetRoleServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) FindTableDataDeptByDeptIDs(ctx context.Context, projectId string, ids map[string]interface{}, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(ids)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	res, err := cli.FindTableDataDeptByDeptIDs(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{
			Data: bts,
		})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateManyTableData(ctx context.Context, projectId, tableName string, closeRequire bool, query, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	if tableName == "" {
		return errors.New("表为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	btsUpdate, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.UpdateMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.MultiUpdateDataRequest{Table: tableName, Query: bts, Data: btsUpdate, CloseRequire: closeRequire})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetWarningFilterIDs(ctx context.Context, projectId, token string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	if token == "" {
		return errors.New("无Token认证信息")
	}
	res, err := cli.GetWarningFilterIDs(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}),
		&api.EmptyRequest{})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryCatalog(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetCatalogServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetCatalog(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetCatalogServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) QueryDept(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetDeptServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetDept(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetDeptServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) QuerySetting(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetSettingServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryApp(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetAppServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetApp(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetAppServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) QuerySystemVariable(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetSystemVariableServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetSystemVariable(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetSystemVariableServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) DeleteSystemVariable(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetSystemVariableServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateSystemVariable(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}

	cli, err := c.CoreClient.GetSystemVariableServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) ReplaceSystemVariable(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.CoreClient.GetSystemVariableServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateSystemVariable(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetSystemVariableServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.Create(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryBackup(ctx context.Context, projectId string, query, result interface{}, count *int64) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetBackupServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	*count = res.GetCount()
	return nil
}

func (c *Client) GetBackup(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetBackupServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) DeleteBackup(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetBackupServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateBackup(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}

	cli, err := c.CoreClient.GetBackupServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) ExportBackup(ctx context.Context, projectId string, query interface{}) (string, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return "", errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetBackupServiceClient()
	if err != nil {
		return "", err
	}
	res, err := cli.Export(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, nil); err != nil {
		return "", err
	}
	id := string(res.GetResult())
	return id, nil
}

func (c *Client) ImportBackup(ctx context.Context, projectId string, query interface{}) (string, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return "", errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetBackupServiceClient()
	if err != nil {
		return "", err
	}
	res, err := cli.Import(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, nil); err != nil {
		return "", err
	}
	id := string(res.GetResult())
	return id, nil
}

func (c *Client) UploadBackup(ctx context.Context, projectId, password string, size int, r io.Reader) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}

	cli, err := c.CoreClient.GetBackupServiceClient()
	if err != nil {
		return err
	}
	stream, err := cli.Upload(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, "password": password}))
	if err != nil {
		return errors.NewResErrorMsg(err, "请求错误")
	}

	//defer stream.CloseAndRecv()

	buffer := make([]byte, 1024)

	bytesReadAll := 0

	for {
		bytesRead, err := r.Read(buffer)
		if err != nil {
			return fmt.Errorf("读取备份文件错误:%s", err.Error())
		}
		err = stream.Send(&core.UploadFileRequest{Data: buffer[:bytesRead]})
		if err != nil {
			return fmt.Errorf("grpc发送备份文件错误:%s", err.Error())
		}
		bytesReadAll += bytesRead

		if bytesReadAll == size {
			err := stream.CloseSend()
			if err != nil {
				return fmt.Errorf("CloseSend错误:%s", err.Error())
			}

			//err := stream.Send(&core.UploadFileRequest{Data: []byte("down")})
			//if err != nil {
			//	return fmt.Errorf("grpc发送结束标志错误:%s", err.Error())
			//}
			//
			//m := new(api.Response)
			//err = stream.RecvMsg(m)
			//if err != nil {
			//	return fmt.Errorf("读取server响应错误:%s", err.Error())
			//}
			//fmt.Printf("上传文件成功，服务器响应结果:%+v\n", m)

			return nil
		}

	}
}

func (c *Client) DownloadBackup(ctx context.Context, projectId, id, password string, w io.Writer) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}

	cli, err := c.CoreClient.GetBackupServiceClient()
	if err != nil {
		return err
	}

	in := new(api.GetOrDeleteRequest)
	in.Id = id
	stream, err := cli.Download(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, "password": password, "id": id}), in)
	if err != nil {
		return errors.NewResErrorMsg(err, "请求错误")
	}

	defer func() {
		_ = stream.CloseSend()
	}()

	for {
		d, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return errors.Wrap(err, "stream.Recv错误")
		}

		logger.Infof("数据长度:%+v", len(d.GetData()))

		n, err := w.Write(d.GetData())
		if err != nil {
			return errors.Wrap(err, " w.Write错误")
		}
		logger.Infof("写入数据长度:%+v", n)
	}

	// 方法二
	//d, err := stream.Recv()
	//if err != nil {
	//	return errors.NewMsg("stream.Recv错误, %s", err.Error())
	//}
	//
	//logger.Infof("数据长度:%+v", len(d.GetData()))
	//
	//n, err := w.Write(d.GetData())
	//if err != nil {
	//	return errors.NewMsg(" w.Write错误, %s", err.Error())
	//}
	//logger.Infof("写入数据长度:%+v", n)
	//return nil

	//for {
	//
	//	err = stream.RecvMsg(&buffer)
	//	if err == io.EOF {
	//		return nil
	//	}
	//	if err != nil {
	//		return err
	//	}
	//
	//	w.Write(buffer)
	//}
}

func (c *Client) FindTagByID(ctx context.Context, projectId, tableId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableId == "" {
		return nil, errors.New("表为空")
	}
	if id == "" {
		return nil, errors.New("设备id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.FindTagByID(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.GetOrDeleteDataRequest{Table: tableId, Id: id})
	return parseRes(err, res, result)
}

func (c *Client) QueryTaskManager(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetTaskManagerServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetTaskManager(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTaskManagerServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) DeleteTaskManager(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTaskManagerServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateTaskManager(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}

	cli, err := c.CoreClient.GetTaskManagerServiceClient()
	if err != nil {
		return err
	}

	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) ReplaceTaskManager(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.CoreClient.GetTaskManagerServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateTaskManager(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetTaskManagerServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.Create(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) FindTableCommandById(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTableSchemaServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.FindCommandByID(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) FindTableDataCommandById(ctx context.Context, projectId, tableId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableId == "" || id == "" {
		return errors.New("表或记录id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.FindCommandByID(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.GetOrDeleteDataRequest{Table: tableId, Id: id})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryTableDataByDB(ctx context.Context, projectId, tableName string, query, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return 0, errors.New("表为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return 0, err
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数为空")
	}
	res, err := cli.QueryByDB(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.QueryDataRequest{
			Table: tableName,
			Query: bts,
		})
	if _, err := parseRes(err, res, result); err != nil {
		return 0, err
	}
	return res.GetCount(), nil
}

func (c *Client) GetTableDataByDB(ctx context.Context, projectId, tableName, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return nil, errors.New("表为空")
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.GetByDB(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.GetOrDeleteDataRequest{Table: tableName, Id: id})
	return parseRes(err, res, result)
}

func (c *Client) DeleteTableDataByDB(ctx context.Context, projectId, tableName, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return errors.New("表为空")
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.DeleteByDB(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.GetOrDeleteDataRequest{Table: tableName, Id: id})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteManyTableDataByDB(ctx context.Context, projectId, tableName string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数为空")
	}
	if tableName == "" {
		return errors.New("表为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.DeleteManyByDB(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.QueryDataRequest{Table: tableName, Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateTableDataByDB(ctx context.Context, projectId, tableName, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return errors.New("表为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.UpdateByDB(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.UpdateDataRequest{Table: tableName, Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) ReplaceTableDataByDB(ctx context.Context, projectId, tableName, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	if tableName == "" {
		return errors.New("表为空")
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.ReplaceByDB(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.UpdateDataRequest{Table: tableName, Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateTableDataByDB(ctx context.Context, projectId, tableName string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return errors.New("表为空")
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.CreateByDB(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.CreateDataRequest{
			Table: tableName,
			Data:  bts,
		})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateManyTableDataByDB(ctx context.Context, projectId, tableName string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if tableName == "" {
		return errors.New("表为空")
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.CreateManyByDB(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.CreateDataRequest{
			Table: tableName,
			Data:  bts,
		})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateManyTableDataByDB(ctx context.Context, projectId, tableName string, updateDataList, result interface{}) error {
	//if projectId == "" {
	//	projectId = config.XRequestProjectDefault
	//}
	//if tableName == "" {
	//	return errors.New("表为空")
	//}
	cli, err := c.CoreClient.GetTableDataServiceClient()
	if err != nil {
		return err
	}
	btsUpdate, err := json.Marshal(updateDataList)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.UpdateManyByDB(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.MultiUpdateDataRequest{Table: tableName, Data: btsUpdate})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateDashboard(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetDashboardServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.Create(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryDashboard(ctx context.Context, projectId string, query, result interface{}) (int, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetDashboardServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return 0, err
	}
	return int(res.GetCount()), nil
}

// UploadFileFromUrl 将远程文件上传到媒体库
//
// sourceUrl 远程文件的下载 url
//
// catalog 上传到媒体库的目录
//
// filename 上传到媒体库后的文件名
//
// 上传成功后返回文件的访问地址,base64字符串,文件大小,错误
func (c *Client) UploadFileFromUrl(ctx context.Context, projectId string, sourceUrl string, catalog string, filename string, action string, addBase64 bool) (string, string, float64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}

	// 布尔值转字符串
	addBase64Str := strconv.FormatBool(addBase64)

	cli, err := c.CoreClient.GetMediaLibraryServiceClient()
	if err != nil {
		return "", "", 0, err
	}

	res, err := cli.UploadFromUrl(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.MediaLibraryUploadFromUrlRequest{
			FileUrl:          sourceUrl,
			MediaLibraryPath: catalog,
			SaveFileName:     filename,
			Action:           action,
			AddBase64:        addBase64Str,
		})
	var result map[string]interface{}
	_, err = parseRes(err, res, &result)
	if err != nil {
		return "", "", 0, err
	}

	fileUrl, ok := result["url"]
	if !ok {
		return "", "", 0, errors.New("上传媒体库成功, 但未返回文件的 url")
	}

	fileUrlStr, ok := fileUrl.(string)
	if !ok {
		return "", "", 0, fmt.Errorf("上传媒体库成功, 但返回文件的 url 不是字符串, %+v", fileUrl)
	}

	sizeRaw, ok := result["size"]
	if !ok {
		return "", "", 0, errors.New("上传媒体库成功, 但未返回文件大小")
	}

	size, ok := sizeRaw.(float64)
	if !ok {
		return "", "", 0, fmt.Errorf("上传媒体库成功, 但返回文件大小不是数字类型, %+v", sizeRaw)
	}

	if addBase64 {
		addBase64StrRaw, ok := result["base64Str"]
		if !ok {
			return "", "", 0, errors.New("上传媒体库成功, 但未返回文件的base64字符串")
		}
		base64Str, ok := addBase64StrRaw.(string)
		if !ok {
			return "", "", 0, fmt.Errorf("上传媒体库成功, 但返回的base64Str不是字符串")
		}
		return fileUrlStr, base64Str, size, nil
	} else {
		return fileUrlStr, "", size, nil
	}
}

func (c *Client) UploadFileFromBase64(ctx context.Context, projectId string, base64Str, mediaLibraryPath, saveFileName, action string) (string, float64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.CoreClient.GetMediaLibraryServiceClient()
	if err != nil {
		return "", 0, err
	}

	res, err := cli.UploadFromBase64(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.MediaLibraryUploadFromBase64Request{
			Base64Str:        base64Str,
			MediaLibraryPath: mediaLibraryPath,
			SaveFileName:     saveFileName,
			Action:           action,
		})

	var result map[string]interface{}
	_, err = parseRes(err, res, &result)
	if err != nil {
		return "", 0, err
	}
	fileUrl, ok := result["url"]
	if !ok {
		return "", 0, errors.New("上传媒体库成功, 但未返回文件的 url")
	}

	fileUrlStr, ok := fileUrl.(string)
	if !ok {
		return "", 0, fmt.Errorf("上传媒体库成功, 但返回文件的 url 不是字符串, %+v", fileUrl)
	}

	sizeRaw, ok := result["size"]
	if !ok {
		return "", 0, errors.New("上传媒体库成功, 但未返回文件大小")
	}

	size, ok := sizeRaw.(float64)
	if !ok {
		return "", 0, fmt.Errorf("上传媒体库成功, 但返回文件大小不是数字类型, %+v", sizeRaw)
	}

	return fileUrlStr, size, nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

var _ io.Reader = (*MultipartReader)(nil)

type MultipartReader struct {
	headers  textproto.MIMEHeader
	boundary string
	body     io.Reader
	header   io.Reader
	tailer   io.Reader
}

func NewMultipart(fieldName, filename string, reader io.Reader) (*MultipartReader, error) {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		return nil, fmt.Errorf("生成 boundary 失败, %+v", err)
	}

	boundary := fmt.Sprintf("%x", buf[:])
	if strings.ContainsAny(boundary, `()<>@,;:\"/[]?= `) {
		boundary = `"` + boundary + `"`
	}

	h := textproto.MIMEHeader{}
	h.Set("Content-Type", "application/octet-stream")
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldName), escapeQuotes(filename)))

	headerReader := bytes.NewBuffer(make([]byte, 0, 1024))
	headerReader.WriteString(fmt.Sprintf("--%s\r\n", boundary))

	for _, k := range slices.Sorted(maps.Keys(h)) {
		for _, v := range h[k] {
			headerReader.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
	}
	headerReader.WriteString("\r\n")

	return &MultipartReader{
		headers:  h,
		boundary: boundary,
		body:     reader,
		header:   headerReader,
		tailer:   bytes.NewReader([]byte(fmt.Sprintf("\r\n--%s--\r\n", boundary))),
	}, nil
}

func (m *MultipartReader) FormDataContentType() string {
	return fmt.Sprintf("multipart/form-data; boundary=%s", m.boundary)
}

func (m *MultipartReader) Read(p []byte) (n int, err error) {
	if m.header != nil {
		n, err = m.header.Read(p)
		if err != nil && err != io.EOF {
			return n, err
		} else if n != 0 {
			return n, nil
		} else if err == io.EOF {
			m.header = nil
		}
	}

	if m.body != nil {
		n, err = m.body.Read(p)
		if err != nil && err != io.EOF {
			return n, err
		} else if n != 0 {
			return n, nil
		} else if err == io.EOF {
			m.body = nil
		}
	}

	if m.tailer != nil {
		n, err = m.tailer.Read(p)
		if err != nil && err != io.EOF {
			return n, err
		} else if n != 0 {
			return n, nil
		} else if err == io.EOF {
			m.tailer = nil
		}
	}

	return 0, io.EOF
}

// UploadFileData 上传文件到媒体库
//
// projectId: 项目ID
// mediaLibraryPath: 媒体库目录
// saveFileName: 保存文件名
// action: 文件重复时的行为处理方式. cover: 覆盖, rename: 文件名自动加1
// reader: 上传文件的 io.Reader
//
// 返回值: 文件访问地址, 错误
func (c *Client) UploadFileData(ctx context.Context, projectId string, mediaLibraryPath, saveFileName, action string, reader io.Reader) (string, error) {
	//body := &bytes.Buffer{}
	//writer := multipart.NewWriter(body)
	//part, err := writer.CreateFormFile("file", saveFileName)
	//if err != nil {
	//	return "", errors.Wrapf(err, "创建 multipart 失改, 文件名 '%s'", saveFileName)
	//}
	//_, err = io.Copy(part, reader)
	//
	//err = writer.Close()
	//if err != nil {
	//	return "", errors.Wrapf(err, "读取上传文件 '%s' 失败", saveFileName)
	//}

	body, err := NewMultipart("file", saveFileName, reader)
	if err != nil {
		return "", fmt.Errorf("创建 multipart 失败, %+v", err)
	}

	req, err := netHttp.NewRequest(netHttp.MethodPost, fmt.Sprintf("/core/mediaLibrary/upload?action=%s&catalog=%s", action, mediaLibraryPath), body)
	if err != nil {
		return "", errors.Wrap(err, "创建 http 请求失败")
	}

	req.Header.Set("Content-Type", body.FormDataContentType())
	resp, err := c.doRestRequest(ctx, projectId, req)
	if err != nil {
		return "", errors.Wrapf(err, "上传文件 '%s' 失败", saveFileName)
	}

	if resp.StatusCode != netHttp.StatusOK {
		err = parseFailedRestResponse(resp)
		return "", fmt.Errorf("读取响应结果失败, %s, %+v", body, err)
	}

	f := MediaFile{}
	err = parseSuccessRestResponse(resp, &f)
	if err != nil {
		return "", fmt.Errorf("解析文件上传结果失败, %+v", err)
	} else if f.Url == "" {
		return "", fmt.Errorf("返回文件访问地址为空")
	}

	return f.Url, nil
}

// UploadFile 上传文件到媒体库
//
// projectId: 项目ID
// mediaLibraryPath: 媒体库目录
// saveFileName: 保存文件名
// action: 文件重复时的行为处理方式. cover: 覆盖, rename: 文件名自动加1
// uploadFile: 本地待上传的文件路径
//
// 返回值: 文件访问地址, 错误
func (c *Client) UploadFile(ctx context.Context, projectId string, mediaLibraryPath, saveFileName, action, uploadFile string) (string, error) {
	file, err := os.Open(uploadFile)
	if os.IsNotExist(err) {
		return "", errors.Errorf("文件 '%s' 不存在", uploadFile)
	} else if err != nil {
		return "", errors.Wrapf(err, "打开文件 '%s' 失败", uploadFile)
	}
	defer file.Close()

	//body := &bytes.Buffer{}
	//writer := multipart.NewWriter(body)
	//part, err := writer.CreateFormFile("file", saveFileName)
	//if err != nil {
	//	return "", errors.Wrapf(err, "创建 multipart 失改, 文件名 '%s'", saveFileName)
	//}
	//_, err = io.Copy(part, file)
	//
	//err = writer.Close()
	//if err != nil {
	//	return "", errors.Wrapf(err, "读取上传文件 '%s' 失败", uploadFile)
	//}

	body, err := NewMultipart("file", saveFileName, file)
	if err != nil {
		return "", fmt.Errorf("创建 multipart 失败, %+v", err)
	}

	req, err := netHttp.NewRequest(netHttp.MethodPost, fmt.Sprintf("/core/mediaLibrary/upload?action=%s&catalog=%s", action, mediaLibraryPath), body)
	if err != nil {
		return "", errors.Wrap(err, "创建 http 请求失败")
	}
	req.Header.Set("Content-Type", body.FormDataContentType())

	resp, err := c.doRestRequest(ctx, projectId, req)
	if err != nil {
		return "", errors.Wrapf(err, "上传文件 '%s' 失败", uploadFile)
	}

	if resp.StatusCode != netHttp.StatusOK {
		err = parseFailedRestResponse(resp)
		return "", fmt.Errorf("读取响应结果失败, %s, %+v", body, err)
	}

	f := MediaFile{}
	err = parseSuccessRestResponse(resp, &f)
	if err != nil {
		return "", fmt.Errorf("解析文件上传结果失败, %+v", err)
	} else if f.Url == "" {
		return "", fmt.Errorf("返回文件访问地址为空")
	}

	return f.Url, nil
}

// DownloadFile 下载媒体库文件到本地
//
// projectId: 项目ID
// path: 媒体库文件路径. /core/fileServer/mediaLibrary/projectId/{filePath} 或 {filePath}
// saveFile: 保存到本地文件的路径
func (c *Client) DownloadFile(ctx context.Context, projectId string, path string, saveFile string) error {
	filePath := bytes.NewBuffer(make([]byte, 0, 128))

	if strings.HasPrefix(path, "/rest") {
		filePath.WriteString(strings.TrimPrefix(path, "/rest"))
	} else if strings.HasPrefix(path, "/core") {
		filePath.WriteString(path)
	} else {
		filePath.WriteString("/core/fileServer/mediaLibrary/")
		filePath.WriteString(projectId)
		if !strings.HasPrefix(path, "/") {
			filePath.WriteString("/")
		}
		filePath.WriteString(path)
	}

	fullPath := filePath.String()
	if strings.HasPrefix(fullPath, "//") {
		fullPath = fullPath[1:]
	}

	req, err := netHttp.NewRequest(netHttp.MethodGet, fullPath, nil)
	if err != nil {
		return errors.Wrap(err, "创建 http 请求失败")
	}

	resp, err := c.doRestRequest(ctx, projectId, req)
	if err != nil {
		err = parseFailedRestResponse(resp)
		return errors.Wrapf(err, "下载文件 '%s' 失败", filePath.String())
	}

	if resp.StatusCode != netHttp.StatusOK {
		err = parseFailedRestResponse(resp)
		return fmt.Errorf("读取下载文件失败, %+v", err)
	}

	f, err := os.OpenFile(saveFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return errors.Wrapf(err, "创建本地文件 '%s' 失败", saveFile)
	}
	defer f.Close()

	bodyReader := resp.Body
	defer bodyReader.Close()
	_, err = io.Copy(f, bodyReader)
	if err != nil {
		return errors.Wrapf(err, "保存文件 '%s' 到本地失败", saveFile)
	}

	return nil
}

// DownloadFileData 下载媒体库文件到本地
//
// projectId: 项目ID
// path: 媒体库文件路径. /core/fileServer/mediaLibrary/projectId/{filePath} 或 {filePath}
//
// 返回值: 文件的字节数组, 错误
func (c *Client) DownloadFileData(ctx context.Context, projectId string, path string) ([]byte, error) {
	filePath := bytes.NewBuffer(make([]byte, 0, 128))

	if strings.HasPrefix(path, "/rest") {
		filePath.WriteString(strings.TrimPrefix(path, "/rest"))
	} else if strings.HasPrefix(path, "/core") {
		filePath.WriteString(path)
	} else {
		filePath.WriteString("/core/fileServer/mediaLibrary/")
		filePath.WriteString(projectId)
		if !strings.HasPrefix(path, "/") {
			filePath.WriteString("/")
		}
		filePath.WriteString(path)
	}

	fullPath := filePath.String()
	if strings.HasPrefix(fullPath, "//") {
		fullPath = fullPath[1:]
	}

	req, err := netHttp.NewRequest(netHttp.MethodGet, fullPath, nil)
	if err != nil {
		return nil, errors.Wrap(err, "创建 http 请求失败")
	}

	resp, err := c.doRestRequest(ctx, projectId, req)
	if err != nil {
		return nil, errors.Wrapf(err, "下载文件 '%s' 失败", fullPath)
	}

	if resp.StatusCode != netHttp.StatusOK {
		err = parseFailedRestResponse(resp)
		return nil, fmt.Errorf("读取下载文件失败, %+v", err)
	}

	writer := bytes.NewBuffer(make([]byte, 0, 1024))

	bodyReader := resp.Body
	defer bodyReader.Close()
	_, err = io.Copy(writer, bodyReader)
	if err != nil {
		return nil, errors.Wrapf(err, "下载文件 '%s' 失败", fullPath)
	}

	return writer.Bytes(), nil
}

func (c *Client) QueryMediaLibrary(ctx context.Context, projectId string, catalog string, isFile, addBase64 bool, query, result interface{}) (int, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetMediaLibraryServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.MediaLibraryQueryRequest{
			Catalog:   catalog,
			IsFile:    isFile,
			AddBase64: addBase64,
			Query:     bts,
		})
	if _, err := parseRes(err, res, result); err != nil {
		return 0, err
	}
	return int(res.GetCount()), nil
}

func (c *Client) MediaLibraryMkdir(ctx context.Context, projectId string, catalog, dirName string) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.CoreClient.GetMediaLibraryServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Mkdir(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.MediaLibraryDirMkDirRequest{
			Catalog: catalog,
			Name:    dirName,
		})
	_, err = parseRes(err, res, nil)
	if err != nil {
		return err
	}
	return nil
}

// CreateMediaLibraryDirSetting 媒体库文件夹设置
func (c *Client) CreateMediaLibraryDirSetting(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetMediaLibraryDirSettingServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.Create(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteMediaLibraryDirSetting(ctx context.Context, projectId, id string, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.CoreClient.GetMediaLibraryDirSettingServiceClient()
	if err != nil {
		return err
	}

	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateMediaLibraryDirSetting(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}

	cli, err := c.CoreClient.GetMediaLibraryDirSettingServiceClient()
	if err != nil {
		return err
	}

	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) ReplaceMediaLibraryDirSetting(ctx context.Context, projectId, id string, updateData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.CoreClient.GetMediaLibraryDirSettingServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据为空")
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetMediaLibraryDirSetting(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.CoreClient.GetMediaLibraryDirSettingServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) GetMediaLibraryDirSettingByPath(ctx context.Context, projectId, path string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if path == "" {
		return nil, errors.New("path为空")
	}
	cli, err := c.CoreClient.GetMediaLibraryDirSettingServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.GetByPath(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&core.MediaLibraryDirSettingQueryByPathRequest{Path: path})
	return parseRes(err, res, result)
}

func (c *Client) QueryMediaLibraryDirSetting(ctx context.Context, projectId string, query, result interface{}) (int, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数为空")
	}
	cli, err := c.CoreClient.GetMediaLibraryDirSettingServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return 0, err
	}
	return int(res.GetCount()), nil
}

func (c *Client) CallAIModel(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.CoreClient.GetUserServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.CallAIModel(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) doRestRequest(ctx context.Context, projectId string, request *netHttp.Request) (*netHttp.Response, error) {
	client, err := c.CoreClient.GetRestClient()
	if err != nil {
		return nil, errors.Wrapf(err, "获取 RestClient 失败")
	}

	token, err := c.GetToken()
	if err != nil {
		return nil, errors.Wrapf(err, "获取平台 Token 失败")
	}

	request.WithContext(ctx)
	request.Header.Set(config.XRequestProject, projectId)
	request.Header.Set(config.XRequestHeaderAuthorization, token)
	resp, err := client.Do(request)
	if err != nil && strings.Contains(err.Error(), "NODE_NOT_FOUND") {
		select {
		case <-ctx.Done():
			return nil, context.DeadlineExceeded
		case <-time.After(time.Second * 3):
			return client.Do(request)
		}
	}
	return resp, err
}

func parseFailedRestResponse(resp *netHttp.Response) error {
	if resp == nil || resp.Body == nil {
		return fmt.Errorf("响应结果为空")
	}

	respBody := resp.Body
	defer respBody.Close()

	bodyLen := 1024
	if resp.ContentLength > 0 {
		bodyLen = int(resp.ContentLength)
	}

	bodyReader := bytes.NewBuffer(make([]byte, 0, bodyLen))
	_, err := io.Copy(bodyReader, respBody)
	if err != nil {
		return errors.Wrapf(err, "读取文件上传结果失败")
	}

	result := map[string]interface{}{}
	respData := bodyReader.Bytes()
	if len(respData) == 0 {
		return fmt.Errorf("响应结果为空")
	}

	body := bodyReader.String()
	err = json.Unmarshal(respData, &result)
	if err != nil {
		return fmt.Errorf("响应结果不是有效对象 '%s'", body)
	}

	var message string
	var detail string

	messageRaw, ok := result["message"]
	if ok {
		if v, ok := messageRaw.(string); ok {
			message = v
		}
	}

	detailRaw, ok := result["detail"]
	if ok {
		if v, ok := detailRaw.(string); ok {
			detail = v
		}
	}

	if message == "" && detail == "" {
		return fmt.Errorf("未知的响应结果, %s", body)
	} else if message != "" && detail == "" {
		return fmt.Errorf("%s", message)
	} else {
		return fmt.Errorf("%s, %s", message, detail)
	}
}

func parseSuccessRestResponse(resp *netHttp.Response, target interface{}) error {
	respBody := resp.Body
	defer respBody.Close()

	bodyLen := 1024
	if resp.ContentLength > 0 {
		bodyLen = int(resp.ContentLength)
	}

	bodyReader := bytes.NewBuffer(make([]byte, 0, bodyLen))
	_, err := io.Copy(bodyReader, respBody)
	if err != nil {
		return fmt.Errorf("读取响应结果失败, %+v", err)
	}

	respData := bodyReader.Bytes()
	if len(respData) == 0 {
		return fmt.Errorf("响应结果为空")
	}

	err = json.Unmarshal(respData, target)
	if err != nil {
		return fmt.Errorf("响应结果不是有效对象, %+v", err)
	}

	return nil
}
