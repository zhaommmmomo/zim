package discovery

import (
	"context"
	"fmt"
	"github.com/zhaommmmomo/zim/common/config"
	"github.com/zhaommmmomo/zim/common/domain"
	"github.com/zhaommmmomo/zim/common/log"
	"github.com/zhaommmmomo/zim/common/utils"
	etcdClient "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type ServiceRegister struct {
	client        *etcdClient.Client
	leaseId       etcdClient.LeaseID
	keepaliveChan <-chan *etcdClient.LeaseKeepAliveResponse
	key           string
	value         string
	ctx           *context.Context
}

func NewServiceRegister(ctx *context.Context, endpoint *domain.Endpoint) (*ServiceRegister, error) {
	// 创建 etcd 客户端
	client, err := etcdClient.New(etcdClient.Config{
		Endpoints:   config.GetEtcdEndpoints(),
		DialTimeout: config.GetEtcdDialTimeout(),
	})
	if err != nil {
		return nil, err
	}

	serviceRegister := &ServiceRegister{
		client: client,
		key:    GenerateRegisterKey(endpoint),
		value:  utils.Marshal(endpoint),
		ctx:    ctx,
	}

	// 设置租约
	err = serviceRegister.lease()

	return serviceRegister, err
}

func (register *ServiceRegister) lease() error {
	// 申请租约
	lease, err := register.client.Grant(*register.ctx, config.GetEtcdLeaseTTL())
	if err != nil {
		return err
	}

	// 添加对应的key信息和租约
	_, err = register.client.Put(*register.ctx, register.key, register.value, etcdClient.WithLease(lease.ID))
	if err != nil {
		return err
	}
	log.Info("ServiceRegister new lease.", zap.Int64("leaseId", int64(register.leaseId)),
		zap.String("key", register.key), zap.String("value", register.value))
	// 保持租约
	keepAliveChan, err := register.client.KeepAlive(*register.ctx, lease.ID)
	if err != nil {
		return err
	}

	register.leaseId = lease.ID
	register.keepaliveChan = keepAliveChan
	return nil
}

func (register *ServiceRegister) UpdateRegisterValue(endpoint *domain.Endpoint) error {
	value := utils.Marshal(endpoint)
	_, err := register.client.Put(*register.ctx, register.key, value, etcdClient.WithLease(register.leaseId))
	if err != nil {
		log.Warn("ServiceRegister update endpoint value fail.", zap.Int64("leaseId", int64(register.leaseId)),
			zap.String("key", register.key), zap.String("value", register.value))
		return err
	}
	log.Info("ServiceRegister update endpoint value success.", zap.Int64("leaseId", int64(register.leaseId)),
		zap.String("key", register.key), zap.String("value", register.value))
	register.value = value
	return nil
}

func (register *ServiceRegister) DelRegisterValue() error {
	_, err := register.client.Delete(*register.ctx, register.key, etcdClient.WithLease(register.leaseId))
	if err != nil {
		log.Warn("ServiceRegister del endpoint fail.", zap.Int64("leaseId", int64(register.leaseId)),
			zap.String("key", register.key))
		return err
	}
	log.Info("ServiceRegister del endpoint success.", zap.Int64("leaseId", int64(register.leaseId)),
		zap.String("key", register.key))
	register.key = ""
	register.value = ""
	register.leaseId = -1
	return nil
}

func GenerateRegisterKey(endpoint *domain.Endpoint) string {
	return fmt.Sprintf("%s/%s:%d", endpoint.Name, endpoint.Ip, endpoint.Port)
}

func SplitRegisterKey(key string) (string, string, int16) {
	i := strings.Index(key, "/")
	j := strings.Index(key, ":")
	port, _ := strconv.Atoi(key[j+1:])
	return key[:i], key[i+1 : j], int16(port)
}
