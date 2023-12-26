package api_client_go

import (
	"context"
	"github.com/air-iot/api-client-go/v4/algorithm"
	"github.com/air-iot/api-client-go/v4/apicontext"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/errors"
	cErrors "github.com/air-iot/errors"
	"github.com/air-iot/json"
	"github.com/air-iot/logger"
	"net/http"
)

// AlgorithmRunById 算法执行
func (c *Client) AlgorithmRunById(ctx context.Context, projectId, id string, data interface{}) ([]byte, error) {

	cli, err := c.AlgorithmClient.GetAlgorithmServiceClient()
	if err != nil {
		return nil, errors.NewMsg("获取客户端错误,%s", err)
	}
	if data == nil {
		return nil, errors.NewMsg("数据为空")
	}

	bts, err := json.Marshal(data)
	if err != nil {
		return nil, errors.NewMsg("数据为空")
	}

	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	res, err := cli.Run(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&algorithm.ClientRunByIdRequest{
			Id:   id,
			Data: bts,
		})
	if err != nil {
		return nil, errors.NewMsg("请求错误, %s", err)
	}
	if res.GetCode() != http.StatusOK {
		logger.Errorf("算法服务grpc响应错误: %+v", res)
		return nil, cErrors.Wrap400Response(err, int(res.GetCode()), "算法服务调用失败, %s", res.GetDetail())
	}

	return res.GetResult(), nil
}
