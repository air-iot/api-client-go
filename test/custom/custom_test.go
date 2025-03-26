package main

import (
	"context"
	"fmt"
	api_client_go "github.com/air-iot/api-client-go/v4"
	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/apicontext"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/local_grpc"
	"github.com/air-iot/api-client-go/v4/metadata"
	"github.com/air-iot/api-client-go/v4/test"
	"log"
	"testing"
	"time"
)

var cli *api_client_go.Client

func TestMain(m *testing.M) {
	log.Println("begin")
	//dsn := "host=airiot.tech user=root password=dell123 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	s := local_grpc.NewServer()
	cli1, clean, err := api_client_go.NewLocalClient(config.Config{
		EtcdConfig: "/airiot/config/dev.json",
		Metadata:   map[string]string{"env": "aliyun"},
		Services:   map[string]config.Service{
			//"spm": {Metadata: map[string]string{"env": "chenpc"}},
			//"core": {Metadata: map[string]string{"env": "local11"}},
			//"flow-engine": {Metadata: map[string]string{"env": "local11"}},
			//"data-service": {Metadata: map[string]string{"env": "local11"}},
			//"flow": {Metadata: map[string]string{"env": "localtest"}},
			//"js-server": {Metadata: map[string]string{"env": "localtest"}},
			//"driver": {Metadata: map[string]string{"env": "local1"}},
		},
		Type: "tenant", // tenant 或 project
		////ProjectId: "default",
		AK: "138dd03b-d3ee-4230-d3d2-520feb580bfe",
		SK: "138dd03b-d3ee-4230-d3d2-520feb580bfd",
		//Timeout: 60,
	}, s)
	if err != nil {
		log.Fatal(err)
	}
	cli = cli1
	//ctx := context.Background()
	//shutdown, err := initProvider()
	//if err != nil {
	// log.Fatal(err)
	//}
	//defer func() {
	//
	//}()
	time.Sleep(time.Second * 2)
	m.Run()
	clean()

	//if err := shutdown(ctx); err != nil {
	// log.Fatal(err)
	//}
	log.Println("end")
}

type TestProvider struct {
	test.UnimplementedFlowTaskServiceServer
}

func (a *TestProvider) Get(ctx context.Context, req *api.GetOrDeleteRequest) (*api.Response, error) {
	info, err := metadata.GetMetaData(ctx)
	if err != nil {
		return nil, err
	}
	if info.ProjectId == "" {
		return nil, fmt.Errorf("无项目信息")
	}

	return &api.Response{
		Status: true,
		Code:   200,
		Info:   "ok",
		Result: []byte("===================================="),
	}, nil
}

func TestGet(t *testing.T) {
	s := local_grpc.NewServer()
	var provider = &TestProvider{}
	test.RegisterFlowTaskServiceServer(s, provider)

	cc := local_grpc.NewClient(s)
	cli := test.NewFlowTaskServiceClient(cc)
	// 重置计时器，开始正式测量

	res, err := cli.Get(apicontext.GetGrpcContext(context.Background(), map[string]string{config.XRequestProject: "aaaaa"}), &api.GetOrDeleteRequest{})
	if err != nil {
		t.Fatalf("client.Get error: %v", err)
	}
	t.Logf("%+v", res)
}

func BenchmarkGet(b *testing.B) {
	s := local_grpc.NewServer()
	var provider = &TestProvider{}
	test.RegisterFlowTaskServiceServer(s, provider)

	cc := local_grpc.NewClient(s, nil)
	cli := test.NewFlowTaskServiceClient(cc)
	// 重置计时器，开始正式测量
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := cli.Get(apicontext.GetGrpcContext(context.Background(), map[string]string{config.XRequestProject: "aaaaa"}), &api.GetOrDeleteRequest{})
		if err != nil {
			b.Fatalf("client.Get error: %v", err)
		}
		//b.Logf("%+v", res)
		_ = res
	}
}
