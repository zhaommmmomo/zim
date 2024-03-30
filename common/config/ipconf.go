package config

import "github.com/spf13/viper"

const (
	ipConfPort = "ip_conf.port"
)

func GetIpConfPort() int {
	return viper.GetInt(ipConfPort)
}
