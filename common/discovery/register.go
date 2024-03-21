package discovery

import (
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/zhaommmmomo/zim/common/config"
	"github.com/zhaommmmomo/zim/common/domain"
	"github.com/zhaommmmomo/zim/common/utils"
	etcdClient "go.etcd.io/etcd/client/v3"
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

func init() {
	config.SetDefault(config.ETCD_LEASE_TTL, 5)
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
	logger.CtxInfof(*register.ctx, "ServiceRegister new lease. leaseId=%d key=%s value=%s", register.leaseId, register.key, register.value)
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
		logger.CtxWarnf(*register.ctx, "ServiceRegister update endpoint value fail. leaseId=%d key=%s value=%s", register.leaseId, register.key, value)
		return err
	}
	logger.CtxInfof(*register.ctx, "ServiceRegister update endpoint value success. leaseId=%d key=%s value=%s", register.leaseId, register.key, value)
	register.value = value
	return nil
}

func (register *ServiceRegister) DelRegisterValue() error {
	_, err := register.client.Delete(*register.ctx, register.key, etcdClient.WithLease(register.leaseId))
	if err != nil {
		logger.CtxWarnf(*register.ctx, "ServiceRegister del endpoint fail. leaseId=%d key=%s", register.leaseId, register.key)
		return err
	}
	logger.CtxInfof(*register.ctx, "ServiceRegister del endpoint success. leaseId=%d key=%s", register.leaseId, register.key)
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
