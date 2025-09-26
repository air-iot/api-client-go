package api_client_go

import (
	"bytes"
	"context"
	"fmt"

	"github.com/air-iot/api-client-go/v4/apicontext"
	"github.com/air-iot/errors"
	"github.com/air-iot/json"

	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/dataservice"
	internalError "github.com/air-iot/api-client-go/v4/errors"
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
		return 0, errors.Wrap(err, "序列化查询参数错误")
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return 0, err
	}
	ctx = apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId})
	res, err := cli.Query(ctx, &api.QueryRequest{Query: bts})
	if _, err := parseRes(err, res, result); err != nil {
		return 0, err
	}
	return res.GetCount(), nil
}

func (c *Client) CreateDataGroups(ctx context.Context, projectId string, createData, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return 0, errors.New("插入数据为空")
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return 0, err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return 0, errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.CreateMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{
			Data: bts,
		})
	if _, err := parseRes(err, res, nil); err != nil {
		return 0, err
	}
	return res.GetCount(), nil
	//if err := json.Unmarshal(res.GetResult(), result); err != nil {
	//	return 0, errors.Wrap(err, "解析请求结果错误")
	//}
}

func (c *Client) ReplaceDataGroup(ctx context.Context, projectId, id string, createData interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据错误")
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{
			Id:   id,
			Data: bts,
		})
	if _, err := parseRes(err, res, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateDataGroup(ctx context.Context, projectId, id string, createData interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据错误")
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{
			Id:   id,
			Data: bts,
		})
	if _, err := parseRes(err, res, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteDataGroup(ctx context.Context, projectId, id string) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{
			Id: id,
		})
	if _, err := parseRes(err, res, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteManyDataGroups(ctx context.Context, projectId string, filter interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(filter)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数错误")
	}
	cli, err := c.DataServiceClient.GetDataGroupServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.DeleteMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{
			Query: bts,
		})
	if _, err := parseRes(err, res, nil); err != nil {
		return 0, err
	}
	return res.GetCount(), nil
}

func (c *Client) QueryDataInterface(ctx context.Context, projectId string, query, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数错误")
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
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

func (c *Client) CreateDataInterfaces(ctx context.Context, projectId string, createData, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return 0, errors.New("插入数据为空")
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
	if err != nil {
		return 0, err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return 0, errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.CreateMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{
			Data: bts,
		})
	if _, err := parseRes(err, res, nil); err != nil {
		return 0, err
	}
	return res.GetCount(), nil
}

func (c *Client) ReplaceDataInterface(ctx context.Context, projectId, id string, createData interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据错误")
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{
			Id:   id,
			Data: bts,
		})
	if _, err := parseRes(err, res, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateDataInterface(ctx context.Context, projectId, id string, createData interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if createData == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(createData)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据错误")
	}
	res, err := cli.Update(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{
			Id:   id,
			Data: bts,
		})
	if _, err := parseRes(err, res, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteDataInterface(ctx context.Context, projectId, id string) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
	if err != nil {
		return err
	}
	res, err := cli.Delete(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.GetOrDeleteRequest{
			Id: id,
		})
	if _, err := parseRes(err, res, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteManyDataInterfaces(ctx context.Context, projectId string, filter interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(filter)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数错误")
	}
	cli, err := c.DataServiceClient.GetDataInterfaceServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.DeleteMany(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{
			Query: bts,
		})
	if _, err := parseRes(err, res, nil); err != nil {
		return 0, err
	}
	return res.GetCount(), nil
}

func (c *Client) DataInterfaceProxy(ctx context.Context, projectId, key string, data map[string]interface{}) (*ProxyResult, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if key == "" {
		return nil, errors.New("key为空")
	}
	if data == nil {
		return nil, errors.New("请求数据为空")
	}
	cli, err := c.DataServiceClient.GetDataServiceClient()
	if err != nil {
		return nil, err
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "序列化请求数据错误")
	}
	res, err := cli.Proxy(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&dataservice.Request{Key: key, Data: bts})
	if err != nil {
		return nil, errors.NewResErrorMsg(err, "请求错误")
	}
	if !res.GetStatus() {
		return nil, errors.Wrap400Response(fmt.Errorf("%s", res.GetDetail()), int(res.GetCode()), "%s", res.GetInfo())
	}
	return &ProxyResult{
		Code:    res.GetHttpCode(),
		Headers: res.GetHeaders(),
		Body:    res.GetResult(),
	}, nil
}

func (c *Client) DatasetViewPreview(ctx context.Context, projectId, mode, datesetId, viewId string, data *QueryParam, result interface{}) ([]byte, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if data == nil {
		return nil, errors.New("请求数据为空")
	}
	cli, err := c.DataServiceClient.GetDatasetViewServiceClient()
	if err != nil {
		return nil, err
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "序列化请求数据错误")
	}
	res, err := cli.Preview(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &dataservice.ViewPreviewReq{
		DatesetId: datesetId,
		ViewId:    viewId,
		Mode:      mode,
		Data:      bts,
	})
	if err != nil {
		return nil, errors.Wrap(err, "请求错误")
	}
	if !res.GetStatus() {
		return nil, internalError.ParseResponse(res)
	}
	if result == nil {
		return res.GetResult(), nil
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return nil, errors.Wrap(err, "解析请求结果错误")
	}
	return res.GetResult(), nil
}

func (c *Client) ReplaceDatasetView(ctx context.Context, projectId, id string, datasetView interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if datasetView == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.DataServiceClient.GetDatasetViewServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(datasetView)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据错误")
	}
	res, err := cli.Replace(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{
			Id:   id,
			Data: bts,
		})
	if err != nil {
		return errors.Wrap(err, "请求错误")
	}
	if !res.GetStatus() {
		return internalError.ParseResponse(res)
	}
	return nil
}

func (c *Client) CreateDatasetView(ctx context.Context, projectId string, datasetView, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if datasetView == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.DataServiceClient.GetDatasetViewServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(datasetView)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.Create(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{
			Data: bts,
		})
	if err != nil {
		return errors.Wrap(err, "请求错误")
	}
	if !res.GetStatus() {
		return internalError.ParseResponse(res)
	}
	//if err := json.Unmarshal(res.GetResult(), result); err != nil {
	//	return 0, errors.Wrap(err, "解析请求结果错误")
	//}
	return nil
}

// 删除全部数据集视图（仅备份还原使用），有特殊认证
func (c *Client) DeleteAllDatasetViews(ctx context.Context, projectId string) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.DataServiceClient.GetDatasetViewServiceClient()
	if err != nil {
		return 0, err
	}
	md := dataservice.ContextDataWPw(ctx)
	md[config.XRequestProject] = projectId
	res, err := cli.DeleteAll(
		apicontext.GetGrpcContext(ctx, md),
		&api.QueryRequest{})
	if err != nil {
		return 0, errors.Wrap(err, "请求错误")
	}
	if !res.GetStatus() {
		return 0, internalError.ParseResponse(res)
	}
	return res.GetCount(), nil
}

func (c *Client) QueryDatasetView(ctx context.Context, projectId string, query, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数错误")
	}
	cli, err := c.DataServiceClient.GetDatasetViewServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if err != nil {
		return 0, errors.Wrap(err, "请求错误")
	}
	if !res.GetStatus() {
		return 0, internalError.ParseResponse(res)
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return 0, errors.Wrap(err, "解析请求结果错误")
	}
	return res.GetCount(), nil
}

// data是数据集配置
func (c *Client) DatasetPreview(ctx context.Context, projectId, mode, datesetId string, data any, result interface{}) ([]byte, error) {
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
		return nil, errors.New("请求数据为空")
	}
	cli, err := c.DataServiceClient.GetDatasetServiceClient()
	if err != nil {
		return nil, err
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "序列化请求数据错误")
	}
	res, err := cli.Preview(apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}), &dataservice.ViewPreviewReq{
		DatesetId: datesetId,
		ViewId:    "",
		Mode:      mode,
		Data:      bts,
	})
	if err != nil {
		return nil, errors.Wrap(err, "请求错误")
	}
	if !res.GetStatus() {
		return nil, internalError.ParseResponse(res)
	}
	if result == nil {
		return res.GetResult(), nil
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return nil, errors.Wrap(err, "解析请求结果错误")
	}
	return res.GetResult(), nil
}

func (c *Client) ReplaceDatasetWhole(ctx context.Context, projectId, id string, datasetWhole interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if datasetWhole == nil {
		return errors.New("更新数据为空")
	}
	cli, err := c.DataServiceClient.GetDatasetServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(datasetWhole)
	if err != nil {
		return errors.Wrap(err, "序列化更新数据错误")
	}
	res, err := cli.ReplaceWhole(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.UpdateRequest{
			Id:   id,
			Data: bts,
		})
	if err != nil {
		return errors.Wrap(err, "请求错误")
	}
	if !res.GetStatus() {
		return internalError.ParseResponse(res)
	}
	return nil
}

func (c *Client) CreateDatasetWhole(ctx context.Context, projectId string, datasetWhole, result interface{}) error {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	if datasetWhole == nil {
		return errors.New("插入数据为空")
	}
	cli, err := c.DataServiceClient.GetDatasetServiceClient()
	if err != nil {
		return err
	}
	bts, err := json.Marshal(datasetWhole)
	if err != nil {
		return errors.Wrap(err, "序列化插入数据错误")
	}
	res, err := cli.CreateWhole(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.CreateRequest{
			Data: bts,
		})
	if err != nil {
		return errors.Wrap(err, "请求错误")
	}
	if !res.GetStatus() {
		return internalError.ParseResponse(res)
	}
	//if err := json.Unmarshal(res.GetResult(), result); err != nil {
	//	return 0, errors.Wrap(err, "解析请求结果错误")
	//}
	return nil
}

// 删除全部数据集（仅备份还原使用），有特殊认证
func (c *Client) DeleteAllDatasets(ctx context.Context, projectId string) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	cli, err := c.DataServiceClient.GetDatasetServiceClient()
	if err != nil {
		return 0, err
	}
	md := dataservice.ContextDataWPw(ctx)
	md[config.XRequestProject] = projectId
	res, err := cli.DeleteAll(
		apicontext.GetGrpcContext(ctx, md),
		&api.QueryRequest{})
	if err != nil {
		return 0, errors.Wrap(err, "请求错误")
	}
	if !res.GetStatus() {
		return 0, internalError.ParseResponse(res)
	}
	return res.GetCount(), nil
}

// 查询数据集整体配置（仅数据集本身），支持过滤和分页
func (c *Client) QueryDataset(ctx context.Context, projectId string, query, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数错误")
	}
	cli, err := c.DataServiceClient.GetDatasetServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.Query(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if err != nil {
		return 0, errors.Wrap(err, "请求错误")
	}
	if !res.GetStatus() {
		return 0, internalError.ParseResponse(res)
	}
	if err := json.Unmarshal(res.GetResult(), result); err != nil {
		return 0, errors.Wrap(err, "解析请求结果错误")
	}
	return res.GetCount(), nil
}

// 查询数据集整体配置（包含数据集本身及其列配置），支持过滤和分页
func (c *Client) QueryDatasetWhole(ctx context.Context, projectId string, query, result interface{}) (int64, error) {
	if projectId == "" {
		projectId = config.XRequestProjectDefault
	}
	bts, err := json.Marshal(query)
	if err != nil {
		return 0, errors.Wrap(err, "序列化查询参数错误")
	}
	cli, err := c.DataServiceClient.GetDatasetServiceClient()
	if err != nil {
		return 0, err
	}
	res, err := cli.QueryWhole(
		apicontext.GetGrpcContext(ctx, map[string]string{config.XRequestProject: projectId}),
		&api.QueryRequest{Query: bts})
	if err != nil {
		return 0, errors.Wrap(err, "请求错误")
	}
	if !res.GetStatus() {
		return 0, internalError.ParseResponse(res)
	}
	d := json.NewDecoder(bytes.NewBuffer(res.GetResult()))
	d.UseNumber()
	if err := d.Decode(&result); err != nil {
		return 0, errors.Wrap(err, "解析请求结果错误")
	}
	return res.GetCount(), nil
}
