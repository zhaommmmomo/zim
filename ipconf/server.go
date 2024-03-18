package ipconf

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/zhaommmmomo/zim/common/config"
)

func Start(path string) {
	config.Init(path)
	s := server.Default(server.WithHostPorts(":" + config.Get("port")))
	s.GET("/ip/list", GetEndpoints)
	s.Spin()
}
