package gateway

import (
	"fmt"
	"github.com/zhaommmmomo/zim/common/base"
	"github.com/zhaommmmomo/zim/common/config"
	"github.com/zhaommmmomo/zim/common/log"
	"net"
	"time"
)

func Start(path string) {
	base.InitBaseComponents(path)

	// 绑定端口并启动服务器
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetGatewayPort()))
	if err != nil {
		panic(err)
	}

	// 初始化 work pool
	//ants.NewPool(1)

	// 初始化 reactor
	_, err = initReactor(&ln)
	if err != nil {
		panic(err)
	}

	log.Info("---------gateway started---------")
	defer ln.Close()
	time.Sleep(time.Second * 300)
	log.Info("---------gateway closed---------")
}
