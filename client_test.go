package api_client_go

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/air-iot/api-client-go/v4/config"
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
		Username:    "root",
		Password:    "dell123",
	})
	if err != nil {
		log.Fatal(err)
	}
	clientEtcd = client

	cli1, clean, err := NewClient(clientEtcd, config.Config{
		Metadata: map[string]string{"env": "aliyun"},
		Services: map[string]config.Service{
			"spm": {Metadata: map[string]string{"env": "local1"}},
		},
		Type:    "tenant",
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
