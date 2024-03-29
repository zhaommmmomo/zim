package config

import "github.com/spf13/viper"

const (
	GatewayPort               = "gateway.port"
	GatewayEpollNum           = "gateway.epoll.num"
	GatewayEpollWaitQueueSize = "gateway.epoll.wait-queue.size"
	GatewayEpollLoadBalancer  = "gateway.epoll.load-balancer"
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

func GetGatewayEpollLoadBalancer() string {
	v := viper.GetString(GatewayEpollLoadBalancer)
	if v == "" {
		v = "roundRobin"
		SetDefault(GatewayEpollLoadBalancer, v)
	}
	return v
}
