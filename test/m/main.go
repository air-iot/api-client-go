package main

import (
	"context"
	"flag"
	"fmt"
	api_client_go "github.com/air-iot/api-client-go/v4"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/json"
)

func main() {
	var (
		coreAddr string
		spmAddr  string
		akType   string
		project  string
		ak       string
		sk       string
		driverId string
		groupId  string
	)

	flag.StringVar(&coreAddr, "core_addr", "", "core grpc地址")
	flag.StringVar(&spmAddr, "spm_addr", "", "spm grpc地址")
	flag.StringVar(&akType, "type", "project", "ak sk类型，tenant 或 project")
	flag.StringVar(&project, "project", "", "项目ID")
	flag.StringVar(&ak, "ak", "", "ak")
	flag.StringVar(&sk, "sk", "", "sk")
	flag.StringVar(&driverId, "driverId", "", "驱动id")
	flag.StringVar(&groupId, "groupId", "", "组id")

	// 解析命令行参数
	flag.Parse()

	cli1, clean, err := api_client_go.NewClient(nil, config.Config{
		LiteMode:   true,
		EtcdConfig: "/airiot/config/dev.json",
		//Metadata:    map[string]string{"env": "aliyun"},
		//Gateway:     coreAddr,
		GatewayGrpc: coreAddr,
		Services: map[string]config.Service{
			//"spm": {Metadata: map[string]string{"env": "chenpc"}},
			"spm": {
				//Metadata: map[string]string{"env": "local"},
				GrpcAddr: spmAddr,
			},
			"core": {
				//Metadata: map[string]string{"env": "local"},
				GrpcAddr: coreAddr,
			},
			//"flow-engine": {Metadata: map[string]string{"env": "local11"}},
			//"data-service": {Metadata: map[string]string{"env": "local11"}},
			//"flow": {Metadata: map[string]string{"env": "localtest"}},
			//"js-server": {Metadata: map[string]string{"env": "localtest"}},
			//"driver": {Metadata: map[string]string{"env": "local1"}},
		},
		Type:      config.KeyType(akType), // tenant 或 project
		ProjectId: project,
		AK:        ak,
		SK:        sk,
		Timeout:   60,
	})
	if err != nil {
		panic(err)
	}
	defer clean()
	var arr []map[string]interface{}
	err = cli1.QueryTableSchemaDeviceByDriverAndGroup(context.Background(), project, driverId, groupId, &arr)
	if err != nil {
		panic(err)
	}
	b, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
