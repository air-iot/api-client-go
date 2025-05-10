package api_client_go

import (
	"context"
	"time"

	"github.com/air-iot/api-client-go/v4/lock"
	"github.com/air-iot/errors"
	"github.com/air-iot/logger"
	gocache "github.com/patrickmn/go-cache"
)

type Service struct {
	lock *lock.KeyLock

	parent        *Client
	ServiceCaches *gocache.Cache
}

func newService(parent *Client) *Service {
	return &Service{
		lock: lock.NewKeyLock(),

		parent:        parent,
		ServiceCaches: gocache.New(5*time.Minute, 10*time.Minute),
	}
}

func (c *Service) IsExist(ctx context.Context, serviceName string) (bool, error) {
	err := c.lock.Lock(ctx, serviceName, serviceName, c.parent.Config.Service.Expire)
	if err != nil {
		return false, errors.Wrap(err, "check service is locked")
	}
	defer func() {
		err := c.lock.Unlock(ctx, serviceName, serviceName, c.parent.Config.Service.Expire)
		if err != nil {
			logger.Errorf("check service unlock err: %v", err)
		}
	}()
	_, ok := c.ServiceCaches.Get(serviceName)
	if ok {
		return true, nil
	}
	var envs []string
	metaEnv, ok := c.parent.Config.Metadata["env"]
	if ok && metaEnv != "" {
		envs = append(envs, metaEnv)
	}
	service, ok := c.parent.Config.Services[serviceName]
	if ok {
		metaEnv, ok = service.Metadata["env"]
		if ok && metaEnv != "" {
			envs = []string{metaEnv}
		}
	}
	// TODO 暂时先处理
	if c.parent.RegistryClient == nil {
		return true, nil
	}
	serviceInstances, err := c.parent.RegistryClient.GetServiceInstances(ctx, serviceName, false, envs...)
	if err != nil {
		return false, err
	}
	if len(serviceInstances) > 0 {
		c.ServiceCaches.Set(serviceName, serviceName, c.parent.Config.Service.Expire)
		return true, nil
	}
	return false, nil
}
