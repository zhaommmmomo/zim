package base

import (
	"fmt"
	"github.com/zhaommmmomo/zim/common/config"
	"github.com/zhaommmmomo/zim/common/log"
)

func InitBaseComponents(configPath string) {
	// 初始化配置文件
	fmt.Println("=========== start init base components ==========")
	config.Init(configPath)
	fmt.Println("=========== init config completed ==========")
	// 初始化日志信息
	log.Init()
	fmt.Println("=========== init log completed ==========")
	fmt.Println("=========== init base components completed ==========")
}
