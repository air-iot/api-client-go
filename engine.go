package api_client_go

import (
	"context"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/metadata"

	"github.com/air-iot/api-client-go/v4/engine"
	"github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/json"
)

func (c *Client) Run(ctx context.Context, projectId, flowConfig string, elementB []byte, variables map[string]interface{}) (result *engine.RunResponse, err error) {
	b, err := json.Marshal(variables)
	if err != nil {
		return nil, errors.NewMsg("序列化变量错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return nil, errors.NewMsg("查询token错误, %s", err)
	}

	cli, err := c.FlowEngineClient.GetDataServiceClient()
	if err != nil {
		return nil, errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Run(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}), &engine.RunRequest{
		ProjectId: projectId,
		Config:    flowConfig,
		Variables: b,
		Element:   elementB,
	})
	if res != nil {
		result = &engine.RunResponse{Job: res.Job}
	}
	if err != nil {
		return result, errors.NewMsg("流程执行错误,%s", err)
	}
	return result, nil
}

func (c *Client) Resume(ctx context.Context, projectId, jobId, elementId string, variables map[string]interface{}) error {
	b, err := json.Marshal(variables)
	if err != nil {
		return errors.NewMsg("序列化变量错误,%s", err)
	}
	cli, err := c.FlowEngineClient.GetDataServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return errors.NewMsg("查询token错误, %s", err)
	}
	if _, err := cli.Resume(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}), &engine.ResumeRequest{
		ProjectId: projectId,
		JobId:     jobId,
		ElementId: elementId,
		Variables: b,
	}); err != nil {
		return errors.NewMsg("流程执行错误,%s", err)
	}
	return nil
}

func (c *Client) Fail(ctx context.Context, projectId, jobId, elementId, errMessage string) error {
	cli, err := c.FlowEngineClient.GetDataServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	token, err := c.Token(projectId)
	if err != nil {
		return errors.NewMsg("查询token错误, %s", err)
	}
	if _, err := cli.Fail(metadata.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId, config.XRequestHeaderAuthorization: token}), &engine.FailRequest{
		ProjectId:    projectId,
		JobId:        jobId,
		ElementId:    elementId,
		ErrorMessage: errMessage,
	}); err != nil {
		return errors.NewMsg("流程执行错误,%s", err)
	}
	return nil
}
