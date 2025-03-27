package api_client_go

import (
	"context"
	"github.com/air-iot/api-client-go/v4/apicontext"
	netHttp "net/http"
	"net/url"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/errors"
	"github.com/air-iot/json"
)

func (c *Client) QueryProject(ctx context.Context, query, result interface{}) error {
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数错误")
	}
	res, err := cli.Query(apicontext.GetGrpcContext(ctx, map[string]string{}), &api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryProjectAvailable(ctx context.Context, result interface{}) error {
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.QueryAvailable(apicontext.GetGrpcContext(ctx, map[string]string{}), &api.EmptyRequest{})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) RestQueryProject(ctx context.Context, query, result interface{}) error {
	u := url.URL{Path: "/spm/project"}
	if query != nil {
		bts, err := json.Marshal(query)
		if err != nil {
			return errors.Wrap(err, "序列化查询参数错误")
		}
		params := url.Values{}
		params.Set("query", string(bts))
		u.RawQuery = params.Encode()
	}
	cli, err := c.SpmClient.GetRestClient()
	if err != nil {
		return err
	}
	if err := cli.Invoke(ctx, netHttp.MethodGet, u.RequestURI(), map[string]interface{}{}, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetProject(ctx context.Context, id string, result interface{}) ([]byte, error) {
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(apicontext.GetGrpcContext(ctx, map[string]string{}), &api.GetOrDeleteRequest{Id: id})
	if err != nil {
		return nil, errors.Wrap(err, "请求错误")
	}
	return parseRes(err, res, result)
}

func (c *Client) DeleteProject(ctx context.Context, id string, result interface{}) error {
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Delete(apicontext.GetGrpcContext(ctx, map[string]string{}), &api.GetOrDeleteRequest{Id: id})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateProject(ctx context.Context, id string, updateData, result interface{}) error {
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据错误")
	}
	res, err := cli.Update(apicontext.GetGrpcContext(ctx, map[string]string{}), &api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateProjectLicense(ctx context.Context, id string, updateData, _ interface{}) error {
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据错误")
	}
	res, err := cli.UpdateLicense(apicontext.GetGrpcContext(ctx, map[string]string{}), &api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) ReplaceProject(ctx context.Context, id string, updateData, result interface{}) error {
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据错误")
	}
	res, err := cli.Replace(apicontext.GetGrpcContext(ctx, map[string]string{}), &api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateProject(ctx context.Context, createData, result interface{}) error {
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.SpmClient.GetProjectServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.Create(apicontext.GetGrpcContext(ctx, map[string]string{}), &api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryPmSetting(ctx context.Context, query, result interface{}) error {
	cli, err := c.SpmClient.GetSettingServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数错误")
	}
	res, err := cli.Query(apicontext.GetGrpcContext(ctx, map[string]string{}), &api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}
