package config

import "github.com/spf13/viper"

const (
	IpConfPort = "ip_conf.port"
)

func GetIpConfPort() int {
	return viper.GetInt(IpConfPort)
}
