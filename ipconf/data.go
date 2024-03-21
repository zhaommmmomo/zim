package ipconf

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/zhaommmmomo/zim/common/discovery"
	"github.com/zhaommmmomo/zim/common/domain"
	"github.com/zhaommmmomo/zim/common/utils"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"sync"
)

type data struct {
	EndpointMap  map[string]State            `json:"endpoint_map"`
	D            *discovery.ServiceDiscovery `json:"-"`
	sync.RWMutex `json:"-"`
}

var endpointData = &data{
	EndpointMap: make(map[string]State),
}

const preKey = "gateway"

func initEndpointData() {
	ctx := context.Background()
	d, _ := discovery.NewServiceDiscover(&ctx)
	go func() {
		defer d.Client.Close()
		endpointData.D = d
		endpoints := d.GetEndpoints(preKey)
		// 初始化数据
		initGateways(&ctx, endpoints)
		startWatch(d)
	}()
}

func initGateways(ctx *context.Context, endpoints []*domain.Endpoint) {
	for _, endpoint := range endpoints {
		key := discovery.GenerateRegisterKey(endpoint)
		state, err := convertState(&endpoint.MetaData)
		if err != nil {
			panic(err)
		}
		endpointData.EndpointMap[key] = *state
	}
	logger.CtxInfof(*ctx, "current endpoint data=%s", utils.Marshal(endpointData))
}

func startWatch(d *discovery.ServiceDiscovery) {
	d.Watch(preKey)
	watchChan := d.WatchChan
	for resp := range watchChan {
		for _, event := range resp.Events {
			switch event.Type {
			case mvccpb.PUT:
				updateEndpointData(d.Ctx, event.Kv)
			case mvccpb.DELETE:
				delEndpointData(d.Ctx, string(event.Kv.Key))
			}
		}
	}
}

func updateEndpointData(ctx *context.Context, kv *mvccpb.KeyValue) {
	endpointData.Lock()
	defer endpointData.Unlock()
	endpoint := &domain.Endpoint{}
	if utils.UnMarshal(kv.Value, endpoint) != nil {
		logger.CtxWarnf(*ctx, "update endpoint data unmarshal endpoint fail. endpoint=%s", string(kv.Value))
		return
	}
	state, err := convertState(&endpoint.MetaData)
	if err != nil {
		logger.CtxWarnf(*ctx, "update endpoint data convertState fail. endpoint=%s", string(kv.Value))
		return
	}
	key := string(kv.Key)
	endpointData.EndpointMap[string(kv.Key)] = *state
	logger.CtxInfof(*ctx, "update endpoint data. key=%s value=%s", key, utils.Marshal(endpointData))
}

func delEndpointData(ctx *context.Context, key string) {
	endpointData.Lock()
	defer endpointData.Unlock()
	delete(endpointData.EndpointMap, key)
	logger.CtxInfof(*ctx, "del endpoint data. key=%s", key)
}
