package api_client_go

import (
	"context"
	"github.com/air-iot/api-client-go/v4/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

var clientEtcd *clientv3.Client

func TestMain(m *testing.M) {
	log.Println("begin")
	//dsn := "host=airiot.tech user=root password=dell123 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"121.89.244.23:2379"},
		//Endpoints:   []string{"localhost:2379"},
		DialTimeout: time.Second * time.Duration(60),
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
		Username:    "root",
		Password:    "dell123",
	})
	if err != nil {
		log.Fatal(err)
	}
	clientEtcd = client
	m.Run()
	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
	log.Println("end")
}

func TestClient_GetTableSchema(t *testing.T) {
	cli, clean, err := NewClient(clientEtcd, config.Config{
		Metadata: map[string]string{"env": "aliyun"},
		Services: nil,
		AK:       "",
		SK:       "",
		Timeout:  60,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer clean()
	for i := 0; i < 2; i++ {
		var obj map[string]interface{}
		if err := cli.GetTableSchema(context.Background(), "625f6dbf5433487131f09ff8", "A模型", &obj); err != nil {
			t.Fatal(err)
		}
		t.Log(obj)
	}

}

func TestClient_QueryProject(t *testing.T) {
	cli, clean, err := NewClient(clientEtcd, config.Config{
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
		t.Fatal(err)
	}
	defer clean()
	for i := 0; i < 1000; i++ {
		var arr []map[string]interface{}
		if err := cli.QueryProject(context.Background(), map[string]interface{}{}, &arr); err != nil {
			t.Error(err)
		}
		t.Log(arr)
		time.Sleep(time.Second * 3)
	}
}

func TestClient_QueryTableSchema(t *testing.T) {
	cli, clean, err := NewClient(clientEtcd, config.Config{
		Metadata: map[string]string{"env": "aliyun"},
		Services: map[string]config.Service{
			"spm":  {Metadata: map[string]string{"env": "local1"}},
			"core": {Metadata: map[string]string{"env": "local1"}},
		},
		Type:    "tenant",
		AK:      "138dd03b-d3ee-4230-d3d2-520feb580bfe",
		SK:      "138dd03b-d3ee-4230-d3d2-520feb580bfd",
		Timeout: 60,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer clean()
	for i := 0; i < 1; i++ {
		var arr []map[string]interface{}
		if err := cli.QueryTableSchema(context.Background(), "625f6dbf5433487131f09ff8", map[string]interface{}{}, &arr); err != nil {
			t.Fatal(err)
		}
		t.Log(arr)
	}
}
