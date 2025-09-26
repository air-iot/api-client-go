package api_client_go

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type watchOptions struct {
	reconnectInterval time.Duration
	allowEmptyEnv     bool
	envs              map[string]struct{}
}

type WatchOption interface {
	Apply(options *watchOptions)
}

type WatchOptionFn func(options *watchOptions)

func (f WatchOptionFn) Apply(options *watchOptions) {
	f(options)
}

func WithEnvs(envs ...string) WatchOptionFn {
	return func(options *watchOptions) {
		options.envs = make(map[string]struct{}, len(envs))
		for i := range envs {
			options.envs[envs[i]] = struct{}{}
		}
	}
}

func WithAllowEmptyEnv(allowEmptyEnv bool) WatchOptionFn {
	return func(options *watchOptions) {
		options.allowEmptyEnv = allowEmptyEnv
	}
}

func WithReconnectInterval(interval time.Duration) WatchOptionFn {
	return func(options *watchOptions) {
		options.reconnectInterval = interval
	}
}

type KratosRegistryWatchClient struct {
	serviceName string
	options     *watchOptions
	err         error
	registry    *etcd.Registry
	instances   []*registry.ServiceInstance
	lock        sync.RWMutex
}

func newKratosRegistryWatchClient(reg *etcd.Registry, serviceName string, options *watchOptions) *KratosRegistryWatchClient {
	return &KratosRegistryWatchClient{
		registry:    reg,
		serviceName: serviceName,
		options:     options,
		instances:   make([]*registry.ServiceInstance, 0, 8),
	}
}

func (rwc *KratosRegistryWatchClient) start(ctx context.Context) error {
	go func() {
		interval := rwc.options.reconnectInterval
		for {
			select {
			case <-ctx.Done():
				return
			default:

			}

			watcher, err := rwc.registry.Watch(ctx, rwc.serviceName)
			if err != nil {
				rwc.err = err
				continue
			}

			for {
				instances, err := watcher.Next()
				if err != nil {
					rwc.err = err
					break
				}

				rwc.lock.Lock()
				rwc.instances = rwc.filter(instances)
				rwc.lock.Unlock()
			}

			time.Sleep(interval)
		}
	}()

	return nil
}

func (rwc *KratosRegistryWatchClient) filter(instances []*registry.ServiceInstance) []*registry.ServiceInstance {
	if len(instances) == 0 {
		return instances
	}

	filtered := make([]*registry.ServiceInstance, 0, len(instances))
	for i := range instances {
		inst := instances[i]
		env, ok := inst.Metadata["env"]
		if (!ok || env == "") && !rwc.options.allowEmptyEnv {
			continue
		}

		if rwc.options.allowEmptyEnv && env == "" {
			filtered = append(filtered, inst)
			continue
		}

		if _, ok := rwc.options.envs[env]; ok {
			filtered = append(filtered, inst)
		}
	}

	return filtered
}

func (rwc *KratosRegistryWatchClient) Error() error {
	return rwc.err
}

func (rwc *KratosRegistryWatchClient) GetServiceInstances() ([]*registry.ServiceInstance, error) {
	rwc.lock.RLock()
	defer rwc.lock.RUnlock()
	return rwc.instances, nil
}

func (rwc *KratosRegistryWatchClient) GetServiceEndpointsByServiceName(schemes ...string) ([]string, error) {
	rwc.lock.RLock()
	defer rwc.lock.RUnlock()
	return getEndpoints(rwc.instances, schemes...), nil
}

type KratosRegistryClient struct {
	registry *etcd.Registry
}

func NewKartosRegistryClient(cli *clientv3.Client, options ...etcd.Option) *KratosRegistryClient {
	return &KratosRegistryClient{registry: etcd.New(cli, options...)}
}

// GetServiceInstances 从注册中心查询指定服务的实例列表. 可以通过 env 过滤服务实例
//
// serviceName: 服务名称
//
// allowEmptyEnv: 是否允许环境为空的实例. 如果为 true, 则返回的实例列表中包含未定义 env 元数据或者 env 为空的实例
//
// envs: 环境列表. 如果为空, 则返回所有实例. 如果不为空, 则返回指定环境的实例
func (reg *KratosRegistryClient) GetServiceInstances(ctx context.Context, serviceName string, allowEmptyEnv bool, envs ...string) ([]*registry.ServiceInstance, error) {
	serviceInstances, err := reg.registry.GetService(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	if len(serviceInstances) == 0 {
		return nil, nil
	} else if len(envs) == 0 {
		return serviceInstances, nil
	}

	finalServiceInstances := make([]*registry.ServiceInstance, 0, len(serviceInstances))
	for i := range serviceInstances {
		inst := serviceInstances[i]
		if env, ok := inst.Metadata["env"]; !ok || env == "" {
			if allowEmptyEnv {
				finalServiceInstances = append(finalServiceInstances, inst)
			}
		} else {
			for ei := range envs {
				if strings.Compare(envs[ei], env) == 0 {
					finalServiceInstances = append(finalServiceInstances, inst)
					break
				}
			}
		}
	}

	return finalServiceInstances, nil
}

// GetServiceEndpoints 获取服务实例的指定协议的端点. 如果 schemes 为空, 则返回所有端点
//
// serviceInstances: 服务实例列表
//
// schemes: 协议列表, 例如: http, https, grpc 等. 如果为空, 则返回所有端点. 如果不为空, 则返回指定协议的端点
func (reg *KratosRegistryClient) GetServiceEndpoints(serviceInstances []*registry.ServiceInstance, schemes ...string) []string {
	return getEndpoints(serviceInstances, schemes...)
}

// GetServiceEndpointsByServiceName 从注册中心查询指定服务的实例列表, 并获取指定协议的端点
//
// serviceName: 服务名称
//
// allowEmptyEnv: 是否允许环境为空的实例. 如果为 true, 则返回的实例列表中包含未定义 env 元数据或者 env 为空的实例
//
// envs: 环境列表. 如果为空, 则返回所有实例. 如果不为空, 则返回指定环境的实例
//
// schemes: 协议列表, 例如: http, https, grpc 等. 如果为空, 则返回所有端点. 如果不为空, 则返回指定协议的端点
func (reg *KratosRegistryClient) GetServiceEndpointsByServiceName(ctx context.Context, serviceName string, allowEmptyEnv bool, envs []string, schemes []string) ([]string, error) {
	serviceInstances, err := reg.GetServiceInstances(ctx, serviceName, allowEmptyEnv, envs...)
	if err != nil {
		return nil, err
	} else if len(serviceInstances) == 0 {
		return nil, nil
	}
	return reg.GetServiceEndpoints(serviceInstances, schemes...), nil
}

func (reg *KratosRegistryClient) Watch(ctx context.Context, serviceName string, options ...WatchOption) (*KratosRegistryWatchClient, error) {
	opts := new(watchOptions)
	for i := range options {
		options[i].Apply(opts)
	}

	if opts.envs == nil {
		opts.envs = make(map[string]struct{})
	}

	if opts.reconnectInterval == 0 {
		opts.reconnectInterval = 5 * time.Second
	}

	watcher := newKratosRegistryWatchClient(reg.registry, serviceName, opts)
	if err := watcher.start(ctx); err != nil {
		return nil, err
	}

	return watcher, nil
}

func getEndpoints(serviceInstances []*registry.ServiceInstance, schemes ...string) []string {
	var endpoints []string
	for i := range serviceInstances {
		inst := serviceInstances[i]
		for j := range inst.Endpoints {
			ep := inst.Endpoints[j]
			if len(schemes) == 0 {
				endpoints = append(endpoints, ep)
			} else {
				for k := range schemes {
					if strings.HasPrefix(ep, schemes[k]) {
						endpoints = append(endpoints, ep)
						break
					}
				}
			}
		}
	}
	return endpoints
}
