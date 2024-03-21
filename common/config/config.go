package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

const (
	IP_CONF_PORT   = "ip_conf.port"
	ETCD_ENDPOINTS = "etcd.endpoints"
	ETCD_TIMEOUT   = "etcd.timeout"
	ETCD_LEASE_TTL = "etcd.lease.ttl"
)

func Init(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("load config {%s} fail: %w", path, err))
	}
	viper.WatchConfig()
}

func Get(key string) string {
	return viper.GetString(key)
}

func GetIpConfPort() int {
	return viper.GetInt(IP_CONF_PORT)
}

func GetEtcdEndpoints() []string {
	return viper.GetStringSlice(ETCD_ENDPOINTS)
}

func GetEtcdDialTimeout() time.Duration {
	return viper.GetDuration(ETCD_TIMEOUT) * time.Second
}

func GetEtcdLeaseTTL() int64 {
	return viper.GetInt64(ETCD_LEASE_TTL)
}

func SetDefault(key string, value any) {
	viper.SetDefault(key, value)
}

func IsDebug() bool {
	return viper.GetString("env") == "dev"
}
