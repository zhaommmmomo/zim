package discovery

import (
	"context"
	"fmt"
	"github.com/zhaommmmomo/zim/common/config"
	"github.com/zhaommmmomo/zim/common/domain"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestNewServiceRegister(t *testing.T) {
	config.Init("/root/go/src/github.com/zhaommmmomo/zim/zim.yaml")
	m := make(map[string]interface{})
	m["key1"] = 123
	endpoint := &domain.Endpoint{
		Name:     "test",
		Ip:       "localhost",
		Port:     8888,
		MetaData: m,
	}
	ctx := context.TODO()
	register, err := NewServiceRegister(&ctx, endpoint)
	if err != nil {
		return
	}
	defer register.client.Close()
	fmt.Println()
	fmt.Printf("leaseId:%d  key:%s  value:%s\n", register.leaseId, register.key, register.value)
	res, err := register.client.Get(ctx, "test", clientv3.WithPrefix())
	if err != nil {
		return
	}
	for _, kv := range res.Kvs {
		fmt.Printf("from kv ==> leaseId:%d key:%s value:%s\n", kv.Lease, string(kv.Key), string(kv.Value))
	}
	fmt.Println()
}

// init
func TestMockRegister1(t *testing.T) {
	config.Init("/root/go/src/github.com/zhaommmmomo/zim/zim.yaml")
	m := make(map[string]interface{})
	m["max_cpu_use"] = 1
	m["cpu_use"] = 11
	endpoint := &domain.Endpoint{
		Name:     "gateway",
		Ip:       "localhost",
		Port:     8888,
		MetaData: m,
	}
	ctx := context.TODO()
	register, err := NewServiceRegister(&ctx, endpoint)
	if err != nil {
		return
	}
	defer register.client.Close()
	time.Sleep(time.Second * 60)
}

// update
func TestMockRegister2(t *testing.T) {
	config.Init("/root/go/src/github.com/zhaommmmomo/zim/zim.yaml")
	m := make(map[string]interface{})
	m["max_cpu_use"] = 2
	m["cpu_use"] = 22
	endpoint := &domain.Endpoint{
		Name:     "gateway",
		Ip:       "localhost",
		Port:     8889,
		MetaData: m,
	}
	ctx := context.TODO()
	register, err := NewServiceRegister(&ctx, endpoint)
	if err != nil {
		return
	}
	defer register.client.Close()
	time.Sleep(time.Second * 3)
	m["max_cpu_use"] = 20
	m["cpu_use"] = 220
	endpoint.MetaData = m
	register.UpdateRegisterValue(endpoint)
	time.Sleep(time.Second * 30)
}

// del
func TestMockRegister3(t *testing.T) {
	config.Init("/root/go/src/github.com/zhaommmmomo/zim/zim.yaml")
	m := make(map[string]interface{})
	m["max_cpu_use"] = 3
	m["cpu_use"] = 33
	endpoint := &domain.Endpoint{
		Name:     "gateway",
		Ip:       "localhost",
		Port:     8890,
		MetaData: m,
	}
	ctx := context.TODO()
	register, err := NewServiceRegister(&ctx, endpoint)
	if err != nil {
		return
	}
	defer register.client.Close()
	time.Sleep(time.Second * 10)
	register.DelRegisterValue()
	time.Sleep(time.Second * 30)
}
