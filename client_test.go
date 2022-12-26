package api_client_go

import (
	"context"
	"github.com/air-iot/api-client-go/v4/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestClient_GetTableSchema(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: time.Second * time.Duration(60),
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
		Username:    "root",
		Password:    "dell123",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	cli, clean, err := NewClient(client, config.Config{
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
