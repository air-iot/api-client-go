package api_client_go

import (
	"fmt"
	"log"

	"dario.cat/mergo"
	"github.com/air-iot/api-client-go/v4/algorithm"
	"github.com/air-iot/api-client-go/v4/auth"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/core"
	"github.com/air-iot/api-client-go/v4/dataservice"
	"github.com/air-iot/api-client-go/v4/driver"
	"github.com/air-iot/api-client-go/v4/engine"
	"github.com/air-iot/api-client-go/v4/flow"
	"github.com/air-iot/api-client-go/v4/live"
	"github.com/air-iot/api-client-go/v4/report"
	"github.com/air-iot/api-client-go/v4/spm"
	"github.com/air-iot/api-client-go/v4/warning"
	"github.com/air-iot/json"
	etcdConfig "github.com/go-kratos/kratos/contrib/config/etcd/v2"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	kratosConfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type Client struct {
	AuthClient        *auth.Client
	SpmClient         *spm.Client
	CoreClient        *core.Client
	FlowClient        *flow.Client
	WarningClient     *warning.Client
	DriverClient      *driver.Client
	DataServiceClient *dataservice.Client
	FlowEngineClient  *engine.Client
	ReportClient      *report.Client
	LiveClient        *live.Client
	AlgorithmClient   *algorithm.Client
}

func NewClient(cli *clientv3.Client, cfg config.Config) (*Client, func(), error) {
	if cfg.EtcdConfig == "" {
		cfg.EtcdConfig = "/airiot/config/pro.json"
	}
	etcdSource, err := etcdConfig.New(cli, etcdConfig.WithPath(cfg.EtcdConfig), etcdConfig.WithPrefix(true))
	if err != nil {
		return nil, nil, fmt.Errorf("查询配置中心错误, %w", err)
	}
	// create a config instance with source
	c2 := kratosConfig.New(kratosConfig.WithSource(
		etcdSource,
		env.NewSource("")),
	)
	defer func() {
		if err := c2.Close(); err != nil {
			log.Println("配置中心关闭错误, ", err.Error())
		}
	}()
	if err := c2.Load(); err != nil {
		return nil, nil, fmt.Errorf("加载配置中心错误, %w", err)
	}
	var c2m map[string]interface{}
	if err := c2.Scan(&c2m); err != nil {
		return nil, nil, fmt.Errorf("配置解析错误, %w", err)
	}
	if err := viper.MergeConfigMap(c2m); err != nil {
		return nil, nil, fmt.Errorf("合并配置错误, %w", err)
	}
	cfgApi := viper.GetStringMap("App.API")
	if cfgApi != nil && len(cfgApi) > 0 {
		var paramMap map[string]interface{}
		if err := json.CopyByJson(&paramMap, cfg); err != nil {
			return nil, nil, fmt.Errorf("复制配置错误, %w", err)
		}
		if err := mergo.Map(&paramMap, cfgApi); err != nil {
			return nil, nil, fmt.Errorf("合并配置错误, %w", err)
		}
		if err := json.CopyByJson(&cfg, paramMap); err != nil {
			return nil, nil, fmt.Errorf("复制结构配置错误, %w", err)
		}
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 120
	}
	r := etcd.New(cli)
	authCli := auth.NewClient(cfg)
	f := func() *auth.Client {
		return authCli
	}
	authCC := auth.NewCustomCredential(f)
	cred := grpc.WithPerRPCCredentials(authCC)
	httpCred := authCC.HttpToken()
	spmClient, cleanSpm, err := spm.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	coreClient, cleanCore, err := core.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	authCli.SetClient(spmClient, coreClient)
	flowClient, cleanFlow, err := flow.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	warningClient, cleanWarning, err := warning.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	driverClient, cleanDriver, err := driver.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	dataServiceClient, cleanDataService, err := dataservice.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	flowEngineClient, cleanFlowEngine, err := engine.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	reportClient, cleanReport, err := report.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	liveClient, cleanLive, err := live.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	algorithmClient, cleanAlgorithm, err := algorithm.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	return &Client{
			AuthClient:        authCli,
			SpmClient:         spmClient,
			CoreClient:        coreClient,
			FlowClient:        flowClient,
			WarningClient:     warningClient,
			DriverClient:      driverClient,
			DataServiceClient: dataServiceClient,
			FlowEngineClient:  flowEngineClient,
			ReportClient:      reportClient,
			LiveClient:        liveClient,
			AlgorithmClient:   algorithmClient,
		}, func() {
			cleanSpm()
			cleanCore()
			cleanFlow()
			cleanWarning()
			cleanDriver()
			cleanDataService()
			cleanFlowEngine()
			cleanReport()
			cleanLive()
			cleanAlgorithm()
		}, nil
}
