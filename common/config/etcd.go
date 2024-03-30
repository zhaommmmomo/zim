package config

import (
	"github.com/spf13/viper"
	"time"
)

const (
	etcdEndpoints = "etcd.endpoints"
	etcdTimeout   = "etcd.timeout"
	etcdLeaseTtl  = "etcd.lease.ttl"
)

func GetEtcdEndpoints() []string {
	return viper.GetStringSlice(etcdEndpoints)
}

func GetEtcdDialTimeout() time.Duration {
	return viper.GetDuration(etcdTimeout) * time.Second
}

func GetEtcdLeaseTTL() int64 {
	v := viper.GetInt64(etcdLeaseTtl)
	if v == 0 {
		v = 5
		SetDefault(etcdLeaseTtl, v)
	}
	return v
}
