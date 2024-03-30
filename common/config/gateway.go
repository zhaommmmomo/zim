package config

import "github.com/spf13/viper"

const (
	gatewayPort               = "gateway.port"
	gatewayEpollNum           = "gateway.epoll.num"
	gatewayEpollWaitQueueSize = "gateway.epoll.wait-queue.size"
	gatewayEpollLoadBalancer  = "gateway.epoll.load-balancer"
	gatewayWorkPoolSize       = "gateway.work-pool.size"
)

func GetGatewayPort() int {
	return viper.GetInt(gatewayPort)
}

func GetGatewayEpollNum() int {
	return viper.GetInt(gatewayEpollNum)
}

func GetGatewayEpollWaitQueueSize() int {
	return viper.GetInt(gatewayEpollWaitQueueSize)
}

func GetGatewayEpollLoadBalancer() int {
	return viper.GetInt(gatewayEpollLoadBalancer)
}

func GetGatewayWorkPoolSize() int {
	return viper.GetInt(gatewayWorkPoolSize)
}
