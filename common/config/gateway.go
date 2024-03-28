package config

import "github.com/spf13/viper"

const (
	GatewayPort               = "gateway.port"
	GatewayEpollNum           = "gateway.epoll.num"
	GatewayEpollWaitQueueSize = "gateway.epoll.wait-queue.size"
)

func GetGatewayPort() int {
	return viper.GetInt(GatewayPort)
}

func GetGatewayEpollNum() int {
	return viper.GetInt(GatewayEpollNum)
}

func GetGatewayEpollWaitQueueSize() int {
	return viper.GetInt(GatewayEpollWaitQueueSize)
}
