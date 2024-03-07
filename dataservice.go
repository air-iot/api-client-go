package api_client_go

import (
	"context"

	"github.com/air-iot/api-client-go/v4/apicontext"
	cErrors "github.com/air-iot/errors"
	"github.com/air-iot/json"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/dataservice"
	"github.com/air-iot/api-client-go/v4/errors"
)

type ProxyResult struct {
	Code    int32  `json:"code"`
	Headers []byte `json:"headers"`
	Body    []byte `json:"body"`
}

type QueryParam struct {
	Fields     []SelectField  `json:"fields,omitempty"` // 包含查询列（原始列originalName或id）、别名和聚合方法
	Where      WhereFilter    `json:"where,omitempty"`  // 查询条件
	Group      []string       `json:"group,omitempty"`  // 聚合的列（原始列的originalName或id）
	Limit      *uint          `json:"limit,omitempty"`  // 查询结果的最大条数
	Offset     *uint          `json:"offset,omitempty"` // 查询结果的偏移量
	Order      []OrderByField `json:"order,omitempty"`  // 排序
	EchartType string         `json:"echartType"`
	NoGroupBy  bool           `json:"noGroupBy"`            // false: groups和stack里的维度字段都会被group by
	Stack      []string       `json:"stack,omitempty"`      // 字段id
	Drill      []string       `json:"drill,omitempty"`      // 下钻，字段id
	GroupAlias []string       `json:"groupAlias,omitempty"` // 聚合字段别名
}

// SelectField @Description	视图查询的列
type SelectField struct {
	Name   string             `json:"name"`             // 数据集的列的列名，使用count函数聚合时不写
	Alias  string             `json:"alias,omitempty"`  // 别名
	Option *SelectFieldOption `json:"option,omitempty"` // 聚合的列需要加上option
}

// SelectFieldOption model info
//
//	@Description	查询选项
//	@Description	如果列涉及到聚合，在选项中配置聚合函
type SelectFieldOption struct {
	Aggregator   string          `json:"aggregator" example:"max"` // 聚合函数
	DefaultValue interface{}     `json:"defaultValue,omitempty"`   // 默认值
	Distinct     bool            `json:"distinct,omitempty"`       // 是否去重
	Filter       *ConditionField `json:"filter,omitempty"`         // having过滤
}

// OrderByField @Description	视图排序
type OrderByField struct {
	Name string `json:"name"`           // 字段名
	Desc bool   `json:"desc,omitempty"` // 是否降序，默认升序
}

// ConditionField
//
//	@Description	视图查询条件，如 name = value
type ConditionField struct {
	Name  string      `json:"name"`  // 字段名
	Value interface{} `json:"value"` // 值
	Op    string      `json:"op"`    // 符号
}

// WhereConditions 只为兼容原配置保留，不再使用
type WhereConditions struct {
	Conditions []ConditionField `json:"conditions"`
	// 不同条件之间的关系，false: and, true: or
	OR bool `json:"or,omitempty"`
}

// @Description	视图所有查询条件组，第一级的条件为或关系，第二级的条件为与关系
type WhereFilter [][]ConditionField

func (wf *WhereFilter) UnmarshalJSON(data []byte) error {
	var condsv2 [][]ConditionField

	err1 := json.Unmarshal(data, &condsv2)
	if err1 == nil {
		*wf = condsv2
		return nil
	}

	var condsV1 WhereConditions
	err := json.Unmarshal(data, &condsV1)
	if err == nil {
		*wf = nil
		return nil
	}

	return err1
}

func (c *Client) QueryDataGroup(ctx context.Context, projectId string, query, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.NewMsg("序列化查询参数为空, %s", err)
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return 0, errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if err != nil {
		return 0, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return 0, cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return 0, errors.NewMsg("解析请求结果错误, %s", err)
	}
	return res.GetCount(), nil
}

func (c *Client) CreateDataGroups(ctx context.Context, projectId string, createData, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return 0, errors.NewMsg("插入数据为空")
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return 0, errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return 0, errors.NewMsg("marshal 插入数据为空")
	}
	res, err := cli.CreateMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{
			Data: bts,
		})
	if err != nil {
		return 0, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return 0, cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	//if err := json.Unmarshal(res.GetResult(), result); err != nil {
	//	return 0, errors.NewMsg("解析请求结果错误, %s", err)
	//}
	return res.GetCount(), nil
}

func (c *Client) ReplaceDataGroup(ctx context.Context, projectId, id string, createData interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.NewMsg("更新数据为空")
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.NewMsg("marshal 更新数据错误,%s", err)
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{
			Id:   id,
			Data: bts,
		})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	return nil
}

func (c *Client) UpdateDataGroup(ctx context.Context, projectId, id string, createData interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.NewMsg("更新数据为空")
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.NewMsg("marshal 更新数据错误,%s", err)
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{
			Id:   id,
			Data: bts,
		})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	return nil
}

func (c *Client) DeleteDataGroup(ctx context.Context, projectId, id string) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{
			Id: id,
		})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	return nil
}

func (c *Client) DeleteManyDataGroups(ctx context.Context, projectId string, filter interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(filter)
	if err != nil {
		return 0, errors.NewMsg("序列化查询参数为空, %s", err)
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return 0, errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.DeleteMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{
			Query: bts,
		})
	if err != nil {
		return 0, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return 0, cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	return res.GetCount(), nil
}

func (c *Client) QueryDataInterface(ctx context.Context, projectId string, query, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.NewMsg("序列化查询参数为空, %s", err)
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
	if err != nil {
		return 0, errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if err != nil {
		return 0, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return 0, cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return 0, errors.NewMsg("解析请求结果错误, %s", err)
	}
	return res.GetCount(), nil
}

func (c *Client) CreateDataInterfaces(ctx context.Context, projectId string, createData, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return 0, errors.NewMsg("插入数据为空")
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
	if err != nil {
		return 0, errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return 0, errors.NewMsg("marshal 插入数据为空")
	}
	res, err := cli.CreateMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{
			Data: bts,
		})
	if err != nil {
		return 0, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return 0, cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	//if err := json.Unmarshal(res.GetResult(), result); err != nil {
	//	return 0, errors.NewMsg("解析请求结果错误, %s", err)
	//}
	return res.GetCount(), nil
}

func (c *Client) ReplaceDataInterface(ctx context.Context, projectId, id string, createData interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.NewMsg("更新数据为空")
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.NewMsg("marshal 更新数据错误,%s", err)
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{
			Id:   id,
			Data: bts,
		})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	return nil
}

func (c *Client) UpdateDataInterface(ctx context.Context, projectId, id string, createData interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.NewMsg("更新数据为空")
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.NewMsg("marshal 更新数据错误,%s", err)
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{
			Id:   id,
			Data: bts,
		})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	return nil
}

func (c *Client) DeleteDataInterface(ctx context.Context, projectId, id string) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
	if err != nil {
		return errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{
			Id: id,
		})
	if err != nil {
		return errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	return nil
}

func (c *Client) DeleteManyDataInterfaces(ctx context.Context, projectId string, filter interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(filter)
	if err != nil {
		return 0, errors.NewMsg("序列化查询参数为空, %s", err)
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
	if err != nil {
		return 0, errors.NewMsg("获取客户端错误,%s", err)
	}
	res, err := cli.DeleteMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{
			Query: bts,
		})
	if err != nil {
		return 0, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return 0, cErrors.Wrap400Response(err, int(res.GetCode()), "响应不成功, %s", res.GetDetail())
	}
	return res.GetCount(), nil
}

func (c *Client) DataInterfaceProxy(ctx context.Context, projectId, key string, data map[string]interface{}) (*ProxyResult, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if key == "" {
		return nil, errors.NewMsg("key为空")
	}
	if data == nil {
		return nil, errors.NewMsg("请求数据为空")
	}
	cli, err := c.DataServiceClient.GetDataServiceClient()
	if err != nil {
		return nil, errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return nil, errors.NewMsg("序列化请求数据错误,%s", err)
	}
	res, err := cli.Proxy(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&dataservice.Request{Key: key, Data: bts})
	if err != nil {
		return nil, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return nil, cErrors.New400Response(int(res.GetCode()), "响应不成功, %s, %s", res.GetInfo(), res.GetDetail())
	}
	return &ProxyResult{
		Code:    res.GetHttpCode(),
		Headers: res.GetHeaders(),
		Body:    res.GetResult(),
	}, nil
	//if res.GetResult() == nil || len(res.GetResult()) == 0 {
	//	return res.GetResult(), nil
	//}
	//if err := json.Unmarshal(res.GetResult(), result); err != nil {
	//	return nil, errors.NewMsg("解析请求结果错误, %s", err)
	//}
	//return res.GetResult(), nil
}

func (c *Client) DatasetViewPreview(ctx context.Context, projectId, mode, datesetId, viewId string, data *QueryParam, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	//if mode == "" {
	//	return nil, errors.NewMsg("mode为空")
	//}
	//if datesetId == "" {
	//	return nil, errors.NewMsg("id为空")
	//}
	if data == nil {
		return nil, errors.NewMsg("请求数据为空")
	}
	cli, err := c.DataServiceClient.GetDatasetViewServiceClient()
	if err != nil {
		return nil, errors.NewMsg("获取客户端错误,%s", err)
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return nil, errors.NewMsg("序列化请求数据错误,%s", err)
	}
	res, err := cli.Preview(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &dataservice.ViewPreviewReq{
		DatesetId: datesetId,
		ViewId:    viewId,
		Mode:      mode,
		Data:      bts,
	})
	if err != nil {
		return nil, errors.NewMsg("请求错误, %s", err)
	}
	if !res.GetStatus() {
		return nil, cErrors.New400Response(int(res.GetCode()), "响应不成功, %s, %s", res.GetInfo(), res.GetDetail())
	}
	if result == nil {
		return res.GetResult(), nil
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return nil, errors.NewMsg("解析请求结果错误, %s", err)
	}
	return res.GetResult(), nil
}
