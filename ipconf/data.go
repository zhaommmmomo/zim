package ipconf

import (
	"context"
	"github.com/zhaommmmomo/zim/common/discovery"
	"github.com/zhaommmmomo/zim/common/domain"
	"github.com/zhaommmmomo/zim/common/log"
	"github.com/zhaommmmomo/zim/common/utils"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.uber.org/zap"
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
		state := State{}
		err := utils.MapToObj(&endpoint.MetaData, &state)
		if err != nil {
			panic(err)
		}
		endpointData.EndpointMap[key] = state
	}
	log.Debug("init gateways endpoint data", zap.String("endpoint", utils.Marshal(endpointData)))
}

func startWatch(d *discovery.ServiceDiscovery) {
	d.Watch(preKey)
	watchChan := d.WatchChan
	for resp := range watchChan {
		for _, event := range resp.Events {
			switch event.Type {
			case mvccpb.PUT:
				updateEndpointData(event.Kv)
			case mvccpb.DELETE:
				delEndpointData(string(event.Kv.Key))
			}
		}
	}
}

func updateEndpointData(kv *mvccpb.KeyValue) {
	endpointData.Lock()
	defer endpointData.Unlock()
	endpoint := &domain.Endpoint{}
	if utils.UnMarshal(kv.Value, endpoint) != nil {
		log.Warn("update endpoint data unmarshal endpoint fail.", zap.String("endpoint", string(kv.Value)))
		return
	}
	state := State{}
	err := utils.MapToObj(&endpoint.MetaData, &state)
	if err != nil {
		log.Warn("update endpoint data convertState fail.", zap.String("endpoint", string(kv.Value)))
		return
	}
	key := string(kv.Key)
	endpointData.EndpointMap[string(kv.Key)] = state
	log.Debug("update endpoint data.", zap.String("key", key), zap.String("value", utils.Marshal(endpointData)))
}

func delEndpointData(key string) {
	endpointData.Lock()
	defer endpointData.Unlock()
	delete(endpointData.EndpointMap, key)
	log.Debug("del endpoint data.", zap.String("key", key))
}
