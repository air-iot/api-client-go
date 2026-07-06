package api_client_go

import (
	"context"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/apicontext"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/engine"
	"github.com/air-iot/errors"
	"github.com/air-iot/json"
)

type TaskMode int32

const (
	USER    TaskMode = 0
	SERVICE TaskMode = 1
)

func (p TaskMode) String() string {
	switch p {
	case USER:
		return "user"
	case SERVICE:
		return "service"
	default:
		return "UNKNOWN"
	}
}

type Params struct {
	ProjectId  string `json:"projectId"`
	FlowId     string `json:"flowId"`
	Job        string `json:"job"`
	ElementId  string `json:"elementId"`
	ElementJob string `json:"elementJob"`
}

type Handler func(param Params, data []byte) (map[string]interface{}, error)

func (c *Client) Run(ctx context.Context, projectId, flowConfig string, elementB []byte, variables map[string]interface{}) (result *engine.RunResponse, err error) {
	b, err := json.Marshal(variables)
	if err != nil {
		return nil, errors.Wrap(err, "序列化变量错误")
	}
	cli, err := c.FlowEngineClient.GetDataServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Run(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &engine.RunRequest{
		ProjectId: projectId,
		Config:    flowConfig,
		Variables: b,
		Element:   elementB,
	})
	if res != nil {
		result = &engine.RunResponse{Job: res.Job}
	}
	if err != nil {
		return result, errors.NewResErrorMsg(err, "流程执行错误")
	}
	return result, nil
}

func (c *Client) Resume(ctx context.Context, projectId, jobId, elementId string, variables map[string]interface{}) error {
	b, err := json.Marshal(variables)
	if err != nil {
		return errors.Wrap(err, "序列化变量错误")
	}
	cli, err := c.FlowEngineClient.GetDataServiceClient()
	if err != nil {
		return err
	}
	if _, err := cli.Resume(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &engine.ResumeRequest{
		ProjectId: projectId,
		JobId:     jobId,
		ElementId: elementId,
		Variables: b,
	}); err != nil {
		return errors.NewResErrorMsg(err, "流程执行错误")
	}
	return nil
}

func (c *Client) Revert(ctx context.Context, projectId, jobId, elementId string, variables map[string]interface{}) error {
	b, err := json.Marshal(variables)
	if err != nil {
		return errors.Wrap(err, "序列化变量错误")
	}
	cli, err := c.FlowEngineClient.GetDataServiceClient()
	if err != nil {
		return err
	}
	if _, err := cli.Revert(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &engine.ResumeRequest{
		ProjectId: projectId,
		JobId:     jobId,
		ElementId: elementId,
		Variables: b,
	}); err != nil {
		return errors.NewResErrorMsg(err, "流程执行错误")
	}
	return nil
}

func (c *Client) Recall(ctx context.Context, projectId, jobId, elementId string, variables map[string]interface{}) error {
	b, err := json.Marshal(variables)
	if err != nil {
		return errors.Wrap(err, "序列化变量错误")
	}
	cli, err := c.FlowEngineClient.GetDataServiceClient()
	if err != nil {
		return err
	}
	if _, err := cli.Recall(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &engine.ResumeRequest{
		ProjectId: projectId,
		JobId:     jobId,
		ElementId: elementId,
		Variables: b,
	}); err != nil {
		return errors.NewResErrorMsg(err, "流程执行错误")
	}
	return nil
}

func (c *Client) Fail(ctx context.Context, projectId, jobId, elementId, errMessage string) error {
	cli, err := c.FlowEngineClient.GetDataServiceClient()
	if err != nil {
		return err
	}
	if _, err := cli.Fail(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &engine.FailRequest{
		ProjectId:    projectId,
		JobId:        jobId,
		ElementId:    elementId,
		ErrorMessage: errMessage,
	}); err != nil {
		return errors.NewResErrorMsg(err, "流程执行错误")
	}
	return nil
}

func (c *Client) QueryFlowJobCron(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数错误")
	}
	cli, err := c.FlowEngineClient.GetFlowJobCronServiceClient()
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

func (c *Client) GetFlowJobCron(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.FlowEngineClient.GetFlowJobCronServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) DeleteFlowJobCron(ctx context.Context, projectId, id string) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.FlowEngineClient.GetFlowJobCronServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	if _, err := parseRes(err, res, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateFlowJobCron(ctx context.Context, projectId, id string, updateData interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}

	cli, err := c.FlowEngineClient.GetFlowJobCronServiceClient()
	if err != nil {
		return err
	}

	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据错误")
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) ReplaceFlowJobCron(ctx context.Context, projectId, id string, updateData interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.FlowEngineClient.GetFlowJobCronServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据错误")
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateFlowJobCron(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.FlowEngineClient.GetFlowJobCronServiceClient()
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

func (c *Client) CreateManyFlowJobCron(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.FlowEngineClient.GetFlowJobCronServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.CreateMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteManyFlowJobCron(ctx context.Context, projectId string, query interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数错误")
	}
	cli, err := c.FlowEngineClient.GetFlowJobCronServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.DeleteMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, nil); err != nil {
		return 0, err
	}
	return res.GetCount(), nil
}

func (c *Client) QueryFlowLogCron(ctx context.Context, projectId string, query, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return errors.Wrap(err, "序列化查询参数错误")
	}
	cli, err := c.FlowEngineClient.GetFlowLogCronServiceClient()
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

func (c *Client) GetFlowLogCron(ctx context.Context, projectId, id string, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return nil, errors.New("id为空")
	}
	cli, err := c.FlowEngineClient.GetFlowLogCronServiceClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.Get(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	return parseRes(err, res, result)
}

func (c *Client) DeleteFlowLogCron(ctx context.Context, projectId, id string) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	cli, err := c.FlowEngineClient.GetFlowLogCronServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{Id: id})
	if _, err := parseRes(err, res, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateFlowLogCron(ctx context.Context, projectId, id string, updateData interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if id == "" {
		return errors.New("id为空")
	}
	if updateData == nil {
		return errors.New("更新数据为空")
	}

	cli, err := c.FlowEngineClient.GetFlowLogCronServiceClient()
	if err != nil {
		return err
	}

	bts, err := json.Marshal(updateData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据错误")
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{Id: id, Data: bts})
	if _, err := parseRes(err, res, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateFlowLogCron(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.FlowEngineClient.GetFlowLogCronServiceClient()
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

func (c *Client) CreateManyFlowLogCron(ctx context.Context, projectId string, createData, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.FlowEngineClient.GetFlowLogCronServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.CreateMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{Data: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteManyFlowLogCron(ctx context.Context, projectId string, query interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数错误")
	}
	cli, err := c.FlowEngineClient.GetFlowLogCronServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.DeleteMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, nil); err != nil {
		return 0, err
	}
	return res.GetCount(), nil
}
