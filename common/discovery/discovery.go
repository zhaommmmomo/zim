package discovery

import (
	"context"
	"github.com/zhaommmmomo/zim/common/config"
	"github.com/zhaommmmomo/zim/common/domain"
	"github.com/zhaommmmomo/zim/common/utils"
	etcdClient "go.etcd.io/etcd/client/v3"
)

type ServiceDiscovery struct {
	Client    *etcdClient.Client
	WatchChan <-chan etcdClient.WatchResponse
	Ctx       *context.Context
}

func NewServiceDiscover(ctx *context.Context) (*ServiceDiscovery, error) {
	client, err := etcdClient.New(etcdClient.Config{
		Endpoints: config.GetEtcdEndpoints(),
	})
	if err != nil {
		panic(err)
	}
	return &ServiceDiscovery{
		Client: client,
		Ctx:    ctx,
	}, err
}

func (d *ServiceDiscovery) GetEndpoints(preKey string) []*domain.Endpoint {
	res, err := d.Client.Get(*d.Ctx, preKey, etcdClient.WithPrefix())
	if err != nil {
		return nil
	}
	var endpoints []*domain.Endpoint
	for _, kv := range res.Kvs {
		endpoint := &domain.Endpoint{}
		if err := utils.UnMarshal(kv.Value, endpoint); err != nil {
			panic(err)
		}
		endpoints = append(endpoints, endpoint)
	}
	return endpoints
}

func (d *ServiceDiscovery) Watch(preKey string) {
	watchChan := d.Client.Watch(*d.Ctx, preKey, etcdClient.WithPrefix())
	d.WatchChan = watchChan
}
