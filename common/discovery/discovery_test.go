package discovery

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func Test(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}
	defer cli.Close()
	lease, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		return
	}
	cli.KeepAlive(context.TODO(), lease.ID)

	go func() {
		var i = 0
		for { //模拟key的变化
			time.Sleep(time.Second)
			cli.Put(context.TODO(), "gateway/127.0.0.1:8888", "a"+fmt.Sprintf("%d", i), clientv3.WithLease(lease.ID))
			i++
			if i > 10 {
				break
			}
		}
	}()

	go func() {
		var i = 0
		for { //模拟key的变化
			time.Sleep(time.Second)
			cli.Put(context.TODO(), "gateway/127.0.0.1:8889", "b"+fmt.Sprintf("%d", i), clientv3.WithLease(lease.ID))
			i++
			if i > 5 {
				break
			}
		}
	}()

	go func() {
		var i = 0
		for { //模拟key的变化
			time.Sleep(time.Second)
			cli.Put(context.TODO(), "gateway1/127.0.0.1:8889", "b"+fmt.Sprintf("%d", i), clientv3.WithLease(lease.ID))
			i++
			if i > 5 {
				break
			}
		}
	}()

	watchChan := cli.Watch(context.TODO(), "gateway/", clientv3.WithPrefix())
	for res := range watchChan {
		for _, event := range res.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Printf("key:%s 修改为:%s Revision:%d\n", event.Kv.Key, event.Kv.Value, event.Kv.ModRevision)
			}
		}
	}

	fmt.Println("over...")
}
