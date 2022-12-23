package metadata

import (
	"context"
	"fmt"

	config "github.com/air-iot/api-client-go/v4/config"
	"google.golang.org/grpc/metadata"
)

type MetaData struct {
	ProjectId string `json:"projectId"`
	Token     string `json:"token"`
}

func GetMetaData(ctx context.Context) (*MetaData, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("无法获取数据")
	}
	res := new(MetaData)
	tokenHeaders := md.Get(config.XRequestHeaderAuthorization)
	if tokenHeaders != nil && len(tokenHeaders) > 0 {
		res.Token = tokenHeaders[0]
		//return nil, fmt.Errorf("无Token认证信息")
	}
	projectIds := md.Get(config.XRequestProject)
	if projectIds != nil && len(projectIds) > 0 {
		res.ProjectId = projectIds[0]
		//return nil, fmt.Errorf("无项目信息")
	}
	return res, nil
}
