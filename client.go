package api_client_go

import (
	"fmt"
	"log"
	"time"

	"dario.cat/mergo"
	"github.com/air-iot/api-client-go/v4/ai"
	"github.com/air-iot/api-client-go/v4/algorithm"
	"github.com/air-iot/api-client-go/v4/api"
	"github.com/air-iot/api-client-go/v4/auth"
	"github.com/air-iot/api-client-go/v4/computerecord"
	"github.com/air-iot/api-client-go/v4/config"
	"github.com/air-iot/api-client-go/v4/core"
	"github.com/air-iot/api-client-go/v4/datarelay"
	"github.com/air-iot/api-client-go/v4/dataservice"
	"github.com/air-iot/api-client-go/v4/driver"
	"github.com/air-iot/api-client-go/v4/engine"
	internalError "github.com/air-iot/api-client-go/v4/errors"
	"github.com/air-iot/api-client-go/v4/flow"
	"github.com/air-iot/api-client-go/v4/jsserver"
	"github.com/air-iot/api-client-go/v4/live"
	"github.com/air-iot/api-client-go/v4/local_grpc"
	"github.com/air-iot/api-client-go/v4/report"
	"github.com/air-iot/api-client-go/v4/spm"
	"github.com/air-iot/api-client-go/v4/sync"
	"github.com/air-iot/api-client-go/v4/syslog"
	"github.com/air-iot/api-client-go/v4/warning"
	"github.com/air-iot/errors"
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
	Config config.Config

	RegistryClient      *KratosRegistryClient
	AuthClient          *auth.Client
	SpmClient           *spm.Client
	CoreClient          *core.Client
	FlowClient          *flow.Client
	WarningClient       *warning.Client
	DriverClient        *driver.Client
	DataServiceClient   *dataservice.Client
	FlowEngineClient    *engine.Client
	ReportClient        *report.Client
	LiveClient          *live.Client
	AlgorithmClient     *algorithm.Client
	DataRelayClient     *datarelay.Client
	JsServerClient      *jsserver.Client
	SyncClient          *sync.Client
	SyslogClient        *syslog.Client
	ComputeRecordClient *computerecord.Client
	AIClient            *ai.Client
	Service             *Service
}

func NewClient(cli *clientv3.Client, cfg config.Config) (*Client, func(), error) {
	return newGrpcClient(cli, cfg)
}

func NewLocalClient(cfg config.Config, ss *local_grpc.Server) (*Client, func(), error) {
	if cfg.EtcdConfig == "" {
		cfg.EtcdConfig = "/airiot/config/pro.json"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 120
	}
	if cfg.ExpirePrecision == 0 {
		cfg.ExpirePrecision = 120
	}
	if cfg.Service.Expire == 0 {
		cfg.Service.Expire = time.Second * 30
	}
	authCli := auth.NewClient(cfg)
	f := func() *auth.Client {
		return authCli
	}

	cc := local_grpc.NewClient(ss, f)
	//authCC := auth.NewCustomCredential(f)
	//cred := grpc.WithPerRPCCredentials(authCC)
	//httpCred := authCC.HttpToken()
	spmClient, cleanSpm, err := spm.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	coreClient, cleanCore, err := core.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	authCli.SetClient(spmClient, coreClient)
	flowClient, cleanFlow, err := flow.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	warningClient, cleanWarning, err := warning.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	driverClient, cleanDriver, err := driver.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	dataServiceClient, cleanDataService, err := dataservice.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	flowEngineClient, cleanFlowEngine, err := engine.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	reportClient, cleanReport, err := report.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	liveClient, cleanLive, err := live.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	algorithmClient, cleanAlgorithm, err := algorithm.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	dataRelayClient, cleanDataRelay, err := datarelay.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	// TODO: jsServerClient, cleanJsServer, err := jsserver.NewLocalClient(cfg, r, cred, httpCred)
	//jsServerClient, cleanJsServer, err := jsserver.NewLocalClient(cfg, r, cred, httpCred)
	//if err != nil {
	//	return nil, nil, err
	//}
	syncClient, cleanSync, err := sync.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	syslogClient, cleanSyslog, err := syslog.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	computeRecordClient, cleanComputeRecord, err := computerecord.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}

	aiClient, cleanAI, err := ai.NewLocalClient(cfg, cc)
	if err != nil {
		return nil, nil, err
	}
	a := &Client{
		Config: cfg,
		//RegistryClient:    NewKartosRegistryClient(cli),
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
		DataRelayClient:   dataRelayClient,
		//JsServerClient:      jsServerClient,
		SyncClient:          syncClient,
		SyslogClient:        syslogClient,
		ComputeRecordClient: computeRecordClient,
		AIClient:            aiClient,
	}
	a.Service = newService(a)
	return a, func() {
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
		cleanDataRelay()
		//cleanJsServer()
		cleanSync()
		cleanSyslog()
		cleanComputeRecord()
		cleanAI()
	}, nil
}

func newGrpcClient(cli *clientv3.Client, cfg config.Config) (*Client, func(), error) {
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
		return nil, nil, errors.Wrap(err, "加载配置中心错误")
	}
	var c2m map[string]interface{}
	if err := c2.Scan(&c2m); err != nil {
		return nil, nil, errors.Wrap(err, "配置解析错误")
	}
	if err := viper.MergeConfigMap(c2m); err != nil {
		return nil, nil, errors.Wrap(err, "合并配置错误")
	}
	cfgApi := viper.GetStringMap("App.API")
	if cfgApi != nil && len(cfgApi) > 0 {
		var paramMap map[string]interface{}
		if err := json.CopyByJson(&paramMap, cfg); err != nil {
			return nil, nil, errors.Wrap(err, "复制配置错误")
		}
		if err := mergo.Map(&paramMap, cfgApi); err != nil {
			return nil, nil, errors.Wrap(err, "合并配置错误")
		}
		if err := json.CopyByJson(&cfg, paramMap); err != nil {
			return nil, nil, errors.Wrap(err, "复制结构配置错误")
		}
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 120
	}

	if cfg.Service.Expire == 0 {
		cfg.Service.Expire = time.Second * 30
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
	dataRelayClient, cleanDataRelay, err := datarelay.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	jsServerClient, cleanJsServer, err := jsserver.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	syncClient, cleanSync, err := sync.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	syslogClient, cleanSyslog, err := syslog.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	computeRecordClient, cleanComputeRecord, err := computerecord.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}

	aiClient, cleanAI, err := ai.NewClient(cfg, r, cred, httpCred)
	if err != nil {
		return nil, nil, err
	}
	a := &Client{
		Config:              cfg,
		RegistryClient:      NewKartosRegistryClient(cli),
		AuthClient:          authCli,
		SpmClient:           spmClient,
		CoreClient:          coreClient,
		FlowClient:          flowClient,
		WarningClient:       warningClient,
		DriverClient:        driverClient,
		DataServiceClient:   dataServiceClient,
		FlowEngineClient:    flowEngineClient,
		ReportClient:        reportClient,
		LiveClient:          liveClient,
		AlgorithmClient:     algorithmClient,
		DataRelayClient:     dataRelayClient,
		JsServerClient:      jsServerClient,
		SyncClient:          syncClient,
		SyslogClient:        syslogClient,
		ComputeRecordClient: computeRecordClient,
		AIClient:            aiClient,
	}
	a.Service = newService(a)
	return a, func() {
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
		cleanDataRelay()
		cleanJsServer()
		cleanSync()
		cleanSyslog()
		cleanComputeRecord()
		cleanAI()
	}, nil
}

func parseRes(err error, res *api.Response, result interface{}) ([]byte, error) {
	if err != nil {
		return nil, errors.NewResErrorMsg(err, "请求错误")
	}
	if !res.GetStatus() {
		return nil, internalError.ParseResponse(res)
	}
	if result != nil && res.GetResult() != nil {
		if err := json.Unmarshal(res.GetResult(), result); err != nil {
			return nil, errors.Wrap(err, "解析请求结果错误")
		}
	}
	return res.GetResult(), nil
}
