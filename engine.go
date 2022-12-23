package api_client_go

import (
	"context"
	"fmt"

	"github.com/air-iot/api-client-go/v4/engine"
	"github.com/air-iot/json"
)

func (c *Client) Run(ctx context.Context, projectId, config string, elementB []byte, variables map[string]interface{}) (result *engine.FlowRunResponse, err error) {
	b, err := json.Marshal(variables)
	if err != nil {
		return nil, fmt.Errorf("序列化变量错误,%s", err)
	}
	//elementB, err := json.Marshal(element)
	//if err != nil {
	//	return nil, fmt.Errorf("json.Marshal element error,%s", err)
	//}

	cli, err := c.FlowEngineClient.GetDataServiceClient()
	if err != nil {
		return nil, fmt.Errorf("获取客户端错误,%s", err)
	}
	res, err := cli.Run(ctx, &engine.RunRequest{
		ProjectId: projectId,
		Config:    config,
		Variables: b,
		Element:   elementB,
	})
	if res != nil {
		result = &engine.FlowRunResponse{Job: res.Job}
	}
	if err != nil {
		return result, fmt.Errorf("流程执行错误,%s", err)
	}
	return result, nil
}

func (c *Client) Resume(ctx context.Context, projectId, jobId, elementId string, variables map[string]interface{}) error {
	b, err := json.Marshal(variables)
	if err != nil {
		return fmt.Errorf("序列化变量错误,%s", err)
	}
	cli, err := c.FlowEngineClient.GetDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	if _, err := cli.Resume(ctx, &engine.ResumeRequest{
		ProjectId: projectId,
		JobId:     jobId,
		ElementId: elementId,
		Variables: b,
	}); err != nil {
		return fmt.Errorf("流程执行错误,%s", err)
	}
	return nil
}

func (c *Client) Fail(ctx context.Context, projectId, jobId, elementId, errMessage string) error {
	cli, err := c.FlowEngineClient.GetDataServiceClient()
	if err != nil {
		return fmt.Errorf("获取客户端错误,%s", err)
	}
	if _, err := cli.Fail(ctx, &engine.FailRequest{
		ProjectId:    projectId,
		JobId:        jobId,
		ElementId:    elementId,
		ErrorMessage: errMessage,
	}); err != nil {
		return fmt.Errorf("流程执行错误,%s", err)
	}
	return nil
}
