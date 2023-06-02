package api_client_go

import (
	"context"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

var clientEtcd *clientv3.Client
var cli *Client

func TestMain(m *testing.M) {
	log.Println("begin")
	//dsn := "host=airiot.tech user=root password=dell123 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: time.Second * time.Duration(60),
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
		Username:    "",
		Password:    "",
	})
	if err != nil {
		log.Fatal(err)
	}
	clientEtcd = client

	cli1, clean, err := NewClient(clientEtcd, config.Config{
		Metadata: map[string]string{"env": "aliyun"},
		Services: map[string]config.Service{
			//"spm":  {Metadata: map[string]string{"env": "local1"}},
			"data-service": {Metadata: map[string]string{"env": "local11"}},
			//"flow-engine": {Metadata: map[string]string{"env": "local1"}},
		},
		Type: "tenant", // tenant 或 project
		//ProjectId: "default",
		AK:      "138dd03b-d3ee-4230-d3d2-520feb580bfe",
		SK:      "138dd03b-d3ee-4230-d3d2-520feb580bfd",
		Timeout: 60,
	})
	if err != nil {
		log.Fatal(err)
	}
	cli = cli1
	m.Run()
	clean()
	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
	log.Println("end")
}

func TestClient_Run(t *testing.T) {
	type Element struct {
		Id     string      `json:"id" bson:"id"`
		Type   string      `json:"type" bson:"type"`
		Config interface{} `json:"_config" bson:"_config" mapstructure:"_config"`
	}

	cfg := `[{"id":"Flow_4FC897A1","type":"flowEnd"},{"_config":{"eventType":"更新记录时","field":[],"query":{"filter":{}},"selectRecord":[],"table":{"id":"新表","title":"表21"}},"id":"6406ab17587f351eaeefd0cc","type":"工作表事件"},{"_config":{"body":{},"name":"数据接口","opKey":"query-all-controlable","opKeyLabel":"查询设备"},"id":"Flow_FAE3D14D","type":"testHandler"},{"id":"Flow_FAE3D14D-6406ab17587f351eaeefd0cc","source":"6406ab17587f351eaeefd0cc","target":"Flow_FAE3D14D","type":"defaultEdge"},{"id":"Flow_4FC897A1-Flow_FAE3D14D","source":"Flow_FAE3D14D","target":"Flow_4FC897A1","type":"defaultEdge"}]`
	elementB, _ := json.Marshal(Element{
		Id:     "6406ab17587f351eaeefd0cc",
		Type:   "worksheetRecord",
		Config: map[string]interface{}{},
	})
	startTimestamp := time.Now().Local().Format(time.RFC3339Nano)
	dStr := `{"#$table":{"_tableName":"table","id":"新表","title":"表21"},"_department":{},"_label":{"name":"a2"},"_settings":{},"_table":"新表","_title":"表21","createTime_default":"2023-03-06T11:07:32.800133+08:00","creator":"admin","creatorName":"admin","disable":false,"extFlowType":"工作表记录修改","extUserMap":{"creator":{"#$user":{"_tableName":"user","id":"admin","name":"admin"}}},"flowTriggerUser":"admin","flowTriggerUserMap":{"#$user":{"_tableName":"user","id":"admin","name":"admin"}},"focus":false,"id":"640558f5b024ee426a4732e5","name":"a2","number-1FF1":1,"off":false,"online":false}`
	var r map[string]interface{}
	json.Unmarshal([]byte(dStr), &r)
	variables1 := map[string]interface{}{
		"#project":                 "625f6dbf5433487131f09ff8",
		"#startTimestamp":          startTimestamp,
		"6406ab17587f351eaeefd0cc": r,
	}
	t.Log(time.Now().Local())
	resp, err := cli.Run(context.Background(), "625f6dbf5433487131f09ff8", cfg, elementB, variables1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(time.Now().Local())
	t.Log(resp)
}

func BenchmarkClient_Run(b *testing.B) {
	type Element struct {
		Id     string      `json:"id" bson:"id"`
		Type   string      `json:"type" bson:"type"`
		Config interface{} `json:"_config" bson:"_config" mapstructure:"_config"`
	}

	cfg := `[{"id":"Flow_4FC897A1","type":"flowEnd"},{"_config":{"eventType":"更新记录时","field":[],"query":{"filter":{}},"selectRecord":[],"table":{"id":"新表","title":"表21"}},"id":"6406ab17587f351eaeefd0cc","type":"工作表事件"},{"_config":{"body":{},"name":"数据接口","opKey":"query-all-controlable","opKeyLabel":"查询设备"},"id":"Flow_FAE3D14D","type":"testHandler"},{"id":"Flow_FAE3D14D-6406ab17587f351eaeefd0cc","source":"6406ab17587f351eaeefd0cc","target":"Flow_FAE3D14D","type":"defaultEdge"},{"id":"Flow_4FC897A1-Flow_FAE3D14D","source":"Flow_FAE3D14D","target":"Flow_4FC897A1","type":"defaultEdge"}]`
	elementB, _ := json.Marshal(Element{
		Id:     "6406ab17587f351eaeefd0cc",
		Type:   "worksheetRecord",
		Config: map[string]interface{}{},
	})
	startTimestamp := time.Now().Local().Format(time.RFC3339Nano)
	dStr := `{"#$table":{"_tableName":"table","id":"新表","title":"表21"},"_department":{},"_label":{"name":"a2"},"_settings":{},"_table":"新表","_title":"表21","createTime_default":"2023-03-06T11:07:32.800133+08:00","creator":"admin","creatorName":"admin","disable":false,"extFlowType":"工作表记录修改","extUserMap":{"creator":{"#$user":{"_tableName":"user","id":"admin","name":"admin"}}},"flowTriggerUser":"admin","flowTriggerUserMap":{"#$user":{"_tableName":"user","id":"admin","name":"admin"}},"focus":false,"id":"640558f5b024ee426a4732e5","name":"a2","number-1FF1":1,"off":false,"online":false}`
	var r map[string]interface{}
	json.Unmarshal([]byte(dStr), &r)
	variables1 := map[string]interface{}{
		"#project":                 "625f6dbf5433487131f09ff8",
		"#startTimestamp":          startTimestamp,
		"6406ab17587f351eaeefd0cc": r,
	}
	for i := 0; i < b.N; i++ {
		resp, err := cli.Run(context.Background(), "625f6dbf5433487131f09ff8", cfg, elementB, variables1)
		if err != nil {
			b.Fatal(err)
		}
		b.Log(resp)
	}
}

func TestClient_GetTableSchema(t *testing.T) {
	for i := 0; i < 2; i++ {
		var obj map[string]interface{}
		if _, err := cli.GetTableSchema(context.Background(), "625f6dbf5433487131f09ff8", "A模型", &obj); err != nil {
			t.Fatal(err)
		}
		t.Log(obj)
	}

}

func TestClient_QueryProject(t *testing.T) {

	//time.Sleep(time.Second * 10)
	for i := 0; i < 3; i++ {
		var arr []map[string]interface{}
		if err := cli.QueryProject(context.Background(), map[string]interface{}{}, &arr); err != nil {
			t.Error(err)
		}
		t.Log(arr)
		//time.Sleep(time.Second * 3)
	}
}

func TestClient_QueryTableSchema(t *testing.T) {
	//time.Sleep(time.Second * 2)
	for i := 0; i < 1; i++ {
		var arr []map[string]interface{}
		if err := cli.QueryTableSchema(context.Background(), "625f6dbf5433487131f09ff8", map[string]interface{}{}, &arr); err != nil {
			//if err := cli.QueryTableSchema(context.Background(), "default", map[string]interface{}{}, &arr); err != nil {
			t.Fatal(err)
		}
		t.Log(arr)
	}
	//time.Sleep(time.Second * 5)
}

func TestClient_GetCurrentUserInfo(t *testing.T) {
	var res map[string]interface{}
	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzU4NjExMDMsImlhdCI6MTY3NTY4ODMwMywibmJmIjoxNjc1Njg4MzAzLCJzdWIiOiI2M2UwYzhiZTM4NGJjYWFjYjAxNzYwZjQiLCJwcm9qZWN0SWQiOiI2MjVmNmRiZjU0MzM0ODcxMzFmMDlmZjgiLCJjdXN0b20iOnsidG9rZW5UeXBlIjoicHJvamVjdCJ9fQ.dk0WNeM1CTXK7J04YIk1cHaZ9xIXKdrVRrqLnUPoDNYJcFLQUdcRWgucDXzuB_uz4-SeaBLJjyjtg_45aFGtng"
	if err := cli.GetCurrentUserInfo(context.Background(), "625f6dbf5433487131f09ff8", token, &res); err != nil {
		t.Fatal(err)
	}
	t.Log("res: ", res)
}

func TestClient_QueryDataGroup(t *testing.T) {
	var res []map[string]interface{}
	type QueryOption struct {
		Limit       *int                   `json:"limit,omitempty"`       // 查询数据长度
		Skip        *int                   `json:"skip,omitempty"`        // 跳过数据长度
		Sort        map[string]int         `json:"sort,omitempty"`        // 排序
		Filter      map[string]interface{} `json:"filter,omitempty"`      // 过滤条件
		WithCount   *bool                  `json:"withCount,omitempty"`   // 是否返回总数
		Project     map[string]interface{} `json:"project,omitempty"`     // 返回字段
		GroupBy     map[string]interface{} `json:"groupBy,omitempty"`     // 聚合分组查询
		GroupFields map[string]interface{} `json:"groupFields,omitempty"` // 聚合分组查询
		WithoutBody *bool                  `json:"withoutBody,omitempty"` // 返回总数,不返回数据
		WithTags    *bool                  `json:"withTags,omitempty"`    // 是否返回最新数据
		Distinct    map[string]interface{} `json:"distinct,omitempty"`    // 是否返回最新数据
		//Joins     []interface{}          `json:"joins,omitempty"`     // 聚合分组查询

		project map[string]interface{}
	}
	t1 := true
	count, err := cli.QueryDataGroup(context.Background(), "625f6dbf5433487131f09ff7", QueryOption{WithCount: &t1}, &res)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("res: ", count, res)
}

func TestClient_QueryPmSetting(t *testing.T) {
	var res map[string]interface{}
	if err := cli.QueryPmSetting(context.Background(), map[string]interface{}{}, &res); err != nil {
		t.Fatal(err)
	}
	t.Log("res: ", res)
}

func TestClient_GetCatalog(t *testing.T) {
	type Item struct {
		FileServer bool `json:"fileServer" example:"true"`
		HTMl       bool `json:"html" example:"true"`
		Mongodb    bool `json:"mongodb" example:"true"`
		Influxdb   bool `json:"influxdb" example:"true"`
	}

	type CatalogSchema struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
		Type string `json:"type,omitempty,omitempty"`
		//User       *UserSchema `json:"user,omitempty"`
		ParentID   string      `json:"parentId,omitempty"`
		Parent     interface{} `json:"parent,omitempty"`
		Order      *float64    `json:"order,omitempty"`
		CreateTime string      `json:"createTime,omitempty"`
		Site       string      `json:"site,omitempty"`
	}

	res := new(CatalogSchema)

	_, err := cli.GetCatalog(context.Background(), "625f6dbf5433487131f09ff7", "634f9eb96f35e9813003b5b5", res)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("res:%+v", res)
}

func TestClient_GetTableRecord(t *testing.T) {

	res := make(map[string]interface{})

	_, err := cli.GetTableRecord(context.Background(), "625f6dbf5433487131f09ff7", "634526abe9ce8a012833a9b9", &res)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("res:%+v", res)
}

func TestClient_QueryBackup(t *testing.T) {

	type Item struct {
		FileServer bool `json:"fileServer" example:"true"`
		HTMl       bool `json:"html" example:"true"`
		Mongodb    bool `json:"mongodb" example:"true"`
		Influxdb   bool `json:"influxdb" example:"true"`
	}

	type BackupSchema struct {
		ID         string `json:"id,omitempty" example:"61a0cc0d05a76adca47efd02"`
		Name       string `json:"name,omitempty" example:"文件名"`
		Status     string `json:"status,omitempty" example:"succeed"`
		Type       string `json:"type,omitempty" example:"export"`
		Log        string `json:"log,omitempty" example:"完成备份！"`
		Path       string `json:"path" example:"/app/backup/610a5205536b84d56e49bfb8/61a0cc0d05a76adca47efd02"` // 日志
		CreateTime string `json:"createTime,omitempty" example:"2021-11-26T19:59:09.393+08:00"`
		UpdateTime string `json:"updateTime,omitempty" example:"2021-11-26T19:59:10.73+08:00"`
		Item       *Item  `json:"item,omitempty"`
	}

	var res []BackupSchema
	type QueryOption struct {
		Limit       *int                   `json:"limit,omitempty"`       // 查询数据长度
		Skip        *int                   `json:"skip,omitempty"`        // 跳过数据长度
		Sort        map[string]int         `json:"sort,omitempty"`        // 排序
		Filter      map[string]interface{} `json:"filter,omitempty"`      // 过滤条件
		WithCount   *bool                  `json:"withCount,omitempty"`   // 是否返回总数
		Project     map[string]interface{} `json:"project,omitempty"`     // 返回字段
		GroupBy     map[string]interface{} `json:"groupBy,omitempty"`     // 聚合分组查询
		GroupFields map[string]interface{} `json:"groupFields,omitempty"` // 聚合分组查询
		WithoutBody *bool                  `json:"withoutBody,omitempty"` // 返回总数,不返回数据
		WithTags    *bool                  `json:"withTags,omitempty"`    // 是否返回最新数据
		Distinct    map[string]interface{} `json:"distinct,omitempty"`    // 是否返回最新数据
		//Joins     []interface{}          `json:"joins,omitempty"`     // 聚合分组查询

		project map[string]interface{}
	}

	type Backup struct {
		ID         string    `json:"id,omitempty" gorm:"column:id;primaryKey;type:string;not null;uniqueIndex;comment:备份唯一标识;"`
		Name       string    `json:"name,omitempty" gorm:"column:name;type:string;not null;comment:备份名称;"`
		Status     string    `json:"status,omitempty" gorm:"column:status;type:string;comment:状态;"`
		Type       string    `json:"type,omitempty" gorm:"column:type;type:string;comment:类型(export/import);"`
		Log        string    `json:"log,omitempty" gorm:"column:log;type:string;comment:日志;"`
		Path       string    `json:"path,omitempty" gorm:"column:path;type:string;comment:路径;"`
		ProjectId  string    `json:"projectId,omitempty" gorm:"column:projectId;type:string;comment:项目ID;"`
		Item       *string   `json:"item,omitempty" gorm:"column:item;type:string;comment:备份项目;"`
		CreateTime time.Time `json:"createTime,omitempty" gorm:"column:createTime;not null;autoCreateTime;comment:创建时间;"`
		UpdateTime time.Time `json:"updateTime,omitempty" gorm:"column:updateTime;autoUpdateTime;comment:修改时间;"`
	}

	withCount := true
	cols := GetColumns(new(Backup))
	query := new(QueryOption)
	query.Project = make(map[string]interface{})
	for _, col := range cols {
		query.Project[col] = 1
	}
	query.WithCount = &withCount

	err := cli.QueryBackup(context.Background(), "625f6dbf5433487131f09ff7", query, &res)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("res:%+v", res)
}

func TestClient_GetBackup(t *testing.T) {

	type Item struct {
		FileServer bool `json:"fileServer" example:"true"`
		HTMl       bool `json:"html" example:"true"`
		Mongodb    bool `json:"mongodb" example:"true"`
		Influxdb   bool `json:"influxdb" example:"true"`
	}

	type BackupSchema struct {
		ID         string `json:"id,omitempty" example:"61a0cc0d05a76adca47efd02"`
		Name       string `json:"name,omitempty" example:"文件名"`
		Status     string `json:"status,omitempty" example:"succeed"`
		Type       string `json:"type,omitempty" example:"export"`
		Log        string `json:"log,omitempty" example:"完成备份！"`
		Path       string `json:"path" example:"/app/backup/610a5205536b84d56e49bfb8/61a0cc0d05a76adca47efd02"` // 日志
		UpdateTime string `json:"updateTime,omitempty" example:"2021-11-26T19:59:10.73+08:00"`
		CreateTime string `json:"createTime,omitempty" example:"2021-11-26T19:59:09.393+08:00"`
		Item       string `json:"item,omitempty"`
		//Item       *Item `json:"item,omitempty"`

	}

	//res := make(map[string]interface{})
	//_, err := cli.GetBackup(context.Background(), "625f6dbf5433487131f09ff7", "6412867fd9a932681abade65", &res)
	//if err != nil {
	//	t.Fatal(err)
	//}

	res := new(BackupSchema)
	_, err := cli.GetBackup(context.Background(), "625f6dbf5433487131f09ff7", "6412867fd9a932681abade65", res)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("res:%+v", res)
}

func TestClient_DeleteBackup(t *testing.T) {

	res := make(map[string]interface{})
	err := cli.DeleteBackup(context.Background(), "625f6dbf5433487131f09ff7", "6412867fd9a932681abade66", &res)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("res:%+v", res)
}

func TestClient_ExportBackup(t *testing.T) {

	type Item struct {
		FileServer bool `json:"fileServer" example:"true"`
		HTMl       bool `json:"html" example:"true"`
		Mongodb    bool `json:"mongodb" example:"true"`
		Influxdb   bool `json:"influxdb" example:"true"`
	}

	type ExportPara struct {
		Name  string                    `json:"name" bson:"name"` // 文件名
		*Item `json:"item" bson:"item"` // 备份项目
	}

	query := new(ExportPara)
	queryStr := `{"name":"2021-12-06","item":{"mongodb":true,"influxdb":false,"fileServer":false,"html":false}}`

	err := json.Unmarshal([]byte(queryStr), query)
	if err != nil {
		t.Fatal(err)
	}

	id, err := cli.ExportBackup(context.Background(), "625f6dbf5433487131f09ff7", query)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("id:%+v", id)
}

func TestClient_ImportBackup(t *testing.T) {

	type Item struct {
		FileServer bool `json:"fileServer" example:"true"`
		HTMl       bool `json:"html" example:"true"`
		Mongodb    bool `json:"mongodb" example:"true"`
		Influxdb   bool `json:"influxdb" example:"true"`
	}

	type ExportPara struct {
		Name  string                    `json:"name" bson:"name"` // 文件名
		*Item `json:"item" bson:"item"` // 备份项目
	}

	query := new(ExportPara)
	queryStr := `{"name":"2021-12-06","id":"6412cd58d4d8742b9175d240","item":{"mongodb":true,"influxdb":false,"fileServer":false,"html":false}}`

	err := json.Unmarshal([]byte(queryStr), query)
	if err != nil {
		t.Fatal(err)
	}

	id, err := cli.ImportBackup(context.Background(), "625f6dbf5433487131f09ff7", query)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("id:%+v", id)
}

func TestClient_DownLoadBackup(t *testing.T) {

	id := "6417f647ea62c9a4582b41d8"

	f, err := os.Create(id + ".zip")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = f.Close()
	}()

	err = cli.DownloadBackup(context.Background(), "625f6dbf5433487131f09ff7", id, "", f)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("id:%+v", id)
}

func TestClient_UploadBackup(t *testing.T) {

	name := "6417f647ea62c9a4582b41d8.zip"
	f, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	id := "6417f647ea62c9a4582b41d8"
	err = cli.UploadBackup(context.Background(), "625f6dbf5433487131f09ff7", "", int(fi.Size()), f)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("id:%+v", id)
}

func TestClient_UploadLicense(t *testing.T) {

	name := "license.txt"
	f, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	id := "6417f647ea62c9a4582b41d8"
	err = cli.UploadLicense(context.Background(), "625f6dbf5433487131f09ff7", int(fi.Size()), f)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("id:%+v", id)
}

// GetColumns 查询实体类 gorm 列名
func GetColumns(a interface{}) []string {
	s := reflect.TypeOf(a)
	if s.Kind() == reflect.Ptr {
		s = s.Elem() //通过反射获取type定义
	}
	result := make([]string, 0)
	//var queryData map[string]interface{}
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		gormTag, ok := field.Tag.Lookup("gorm")
		if !ok {
			continue
		}
		tags := strings.Split(gormTag, ";")
		for _, tag := range tags {
			if strings.HasPrefix(tag, "column:") {
				columns := strings.Split(tag, ":")
				if len(columns) == 2 && columns[1] != "" {
					result = append(result, columns[1])
					break
				}
			}
		}
	}
	return result
}

func TestClient_clean(t *testing.T) {
	type e struct {
		id string `json:"id"`
	}
	ctx := context.Background()
	var projects []e
	if err := cli.QueryProject(ctx, map[string]interface{}{}, &projects); err != nil {
		t.Fatal(err)
	}
	for _, p := range projects {

		projectId := p.id
		if projectId == "63e1d40cc3879495dfe8b5e4" {
			continue
		}
		var flows []e
		if _, err := cli.QueryFlow(ctx, projectId, map[string]interface{}{}, &flows); err != nil {
			for _, f := range flows {
				var res map[string]interface{}
				err := cli.UpdateFlow(context.Background(), projectId, f.id, map[string]interface{}{"disable": true}, f)
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("id:%+v", res)
			}

		}
	}

}

func TestClient_ReplaceDataGroups(t *testing.T) {

	if err := cli.ReplaceDataGroup(context.Background(), "625f6dbf5433487131f09ff9", "6461d962693d5e41ef126b8e", map[string]interface{}{
		"name": "测试3",
		"type": "http",
	}); err != nil {
		t.Error(err)
	}

}

func Test_DeleteManyDataGroups(t *testing.T) {
	if _, err := cli.DeleteManyDataGroups(context.Background(), "625f6dbf5433487131f09ff9", map[string]interface{}{
		"filter": map[string]interface{}{"id": []interface{}{"6461d962693d5e41ef126b8e", "6461d613693d5e41ef126b8d"}},
	}); err != nil {
		t.Error(err)
	}

}

func TestClient_ReplaceDataInterfaces(t *testing.T) {

	if err := cli.ReplaceDataInterface(context.Background(), "625f6dbf5433487131f09ff9", "6461dc43693d5e41ef126b8f", map[string]interface{}{
		"dataGroup": map[string]interface{}{"id": "6461d962693d5e41ef126b8e", "name": "测试3", "type": "http", "createTime": "2023-05-15T15:04:02+08:00"},
		"key":       "test11", "name": "test1121", "setting": map[string]interface{}{"method": "GET"}}); err != nil {
		t.Error(err)
	}

}

func Test_DeleteManyDataInterfaces(t *testing.T) {
	if _, err := cli.DeleteManyDataInterfaces(context.Background(), "625f6dbf5433487131f09ff9", map[string]interface{}{
		"filter": map[string]interface{}{"id": []interface{}{"6461dc43693d5e41ef126b8f"}},
	}); err != nil {
		t.Error(err)
	}

}

func TestClient_RtspPull(t *testing.T) {

	// RTSPTransSrv RTSP 转换服务 struct
	type RTSPTransSrv struct {
		URL        string `json:"url"`
		Resolution string `json:"resolution"`
		FrameRate  int    `json:"frameRate"`
	}
	bs := RTSPTransSrv{
		URL: "rtsp://admin:aaa123456@39.103.179.147:19030/h264/ch1/sub/av_stream",
	}

	if res, err := cli.RtspPull(context.Background(), "6273677fe5bd0d8ebc7d6a5f", bs); err != nil {
		t.Error(err)
	} else {
		t.Log(res)
	}
}
