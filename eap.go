package api_client_go

import (
	"context"

	"github.com/air-iot/api-client-go/v4/apicontext"
	"github.com/air-iot/api-client-go/v4/config"
	pb "github.com/air-iot/api-client-go/v4/eap"
)

type EAPTaskAttachment struct {
	URL         string
	Filename    string
	ContentType string
}
type SendEAPTaskMessageRequest struct {
	TaskID       string
	Content      string
	SystemPrompt string
	Attachments  []EAPTaskAttachment
}
type SendEAPTaskMessageResult struct {
	TaskID    string
	RunID     string
	MessageID string
}

func (c *Client) SendEAPTaskMessage(ctx context.Context, projectID string, request *SendEAPTaskMessageRequest) (*SendEAPTaskMessageResult, error) {
	client, err := c.EAPClient.GetTaskServiceClient()
	if err != nil {
		return nil, err
	}
	attachments := make([]*pb.Attachment, 0, len(request.Attachments))
	for _, attachment := range request.Attachments {
		attachments = append(attachments, &pb.Attachment{Url: attachment.URL, Filename: attachment.Filename, ContentType: attachment.ContentType})
	}
	if projectID == "" {
		projectID = config.XRequestProjectDefault
	}
	response, err := client.SendMessage(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectID}), &pb.SendTaskMessageRequest{TaskId: request.TaskID, Content: request.Content, SystemPrompt: request.SystemPrompt, Attachments: attachments})
	if err != nil {
		return nil, err
	}
	return &SendEAPTaskMessageResult{TaskID: response.TaskId, RunID: response.RunId, MessageID: response.MessageId}, nil
}
