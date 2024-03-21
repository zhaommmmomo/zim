package ipconf

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/zhaommmmomo/zim/common/config"
)

func Start(path string) {
	config.Init(path)
	initEndpointData()
	s := server.Default(server.WithHostPorts(fmt.Sprintf(":%d", config.GetIpConfPort())))
	s.GET("/ip/list", getGateWayInfo)
	s.Spin()
}
