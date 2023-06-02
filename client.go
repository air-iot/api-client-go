package api_client_go

import (
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
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type Client struct {
	SpmClient         *spm.Client
	CoreClient        *core.Client
	FlowClient        *flow.Client
	WarningClient     *warning.Client
	DriverClient      *driver.Client
	DataServiceClient *dataservice.Client
	FlowEngineClient  *engine.Client
	ReportClient      *report.Client
	LiveClient        *live.Client
}

func NewClient(cli *clientv3.Client, cfg config.Config) (*Client, func(), error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 120
	}
	r := etcd.New(cli)
	authCli := auth.NewClient(cfg)
	f := func() *auth.Client {
		return authCli
	}
	cred := grpc.WithPerRPCCredentials(auth.NewCustomCredential(f))
	spmClient, cleanSpm, err := spm.NewClient(cfg, r, cred)
	if err != nil {
		return nil, nil, err
	}
	coreClient, cleanCore, err := core.NewClient(cfg, r, cred)
	if err != nil {
		return nil, nil, err
	}
	authCli.SetClient(spmClient, coreClient)

	flowClient, cleanFlow, err := flow.NewClient(cfg, r, cred)
	if err != nil {
		return nil, nil, err
	}
	warningClient, cleanWarning, err := warning.NewClient(cfg, r, cred)
	if err != nil {
		return nil, nil, err
	}
	driverClient, cleanDriver, err := driver.NewClient(cfg, r, cred)
	if err != nil {
		return nil, nil, err
	}
	dataServiceClient, cleanDataService, err := dataservice.NewClient(cfg, r, cred)
	if err != nil {
		return nil, nil, err
	}
	flowEngineClient, cleanFlowEngine, err := engine.NewClient(cfg, r, cred)
	reportClient, cleanReport, err := report.NewClient(cfg, r, cred)
	if err != nil {
		return nil, nil, err
	}
	liveClient, cleanLive, err := live.NewClient(cfg, r, cred)
	if err != nil {
		return nil, nil, err
	}
	return &Client{
			SpmClient:         spmClient,
			CoreClient:        coreClient,
			FlowClient:        flowClient,
			WarningClient:     warningClient,
			DriverClient:      driverClient,
			DataServiceClient: dataServiceClient,
			FlowEngineClient:  flowEngineClient,
			ReportClient:      reportClient,
			LiveClient:        liveClient,
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
		}, nil
}
