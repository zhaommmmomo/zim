package config

import (
	"github.com/spf13/viper"
	"time"
)

const (
	EtcdEndpoints = "etcd.endpoints"
	EtcdTimeout   = "etcd.timeout"
	EtcdLeaseTtl  = "etcd.lease.ttl"
)

func GetEtcdEndpoints() []string {
	return viper.GetStringSlice(EtcdEndpoints)
}

func GetEtcdDialTimeout() time.Duration {
	return viper.GetDuration(EtcdTimeout) * time.Second
}

func GetEtcdLeaseTTL() int64 {
	v := viper.GetInt64(EtcdLeaseTtl)
	if v == 0 {
		v = 5
		SetDefault(EtcdLeaseTtl, v)
	}
	return v
}
