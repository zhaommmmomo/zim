package ipconf

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zhaommmmomo/zim/common/discovery"
	"github.com/zhaommmmomo/zim/common/domain"
	"sort"
)

type GateWayInfo struct {
	Name  string `json:"name"`
	Ip    string `json:"ip"`
	Port  int16  `json:"port"`
	Score int16  `json:"score"`
}

func getGateWayInfo(context context.Context, appCtx *app.RequestContext) {
	defer func() {
		if err := recover(); err != nil {
			appCtx.JSON(consts.StatusBadRequest, utils.H{"err": err})
		}
	}()
	// 构建请求上下文参数
	_ = domain.BuildIpConfContext(&context, appCtx)
	// 获取集群中对应的endpoints  按照对应的分数进行排序了
	gateways := dispatch()
	// 返回数据
	appCtx.JSON(consts.StatusOK, domain.PackRes(top5GateWays(gateways)))
}

func dispatch() []*GateWayInfo {
	// 1. 计算各个服务器的分值
	gateways := calculateGatewayScore()
	// 2. 对服务器进行排序（分数由高到低）
	sort.Slice(gateways, func(i, j int) bool {
		return gateways[i].Score > gateways[j].Score
	})
	return gateways
}

func calculateGatewayScore() []*GateWayInfo {
	endpointData.RLock()
	defer endpointData.RUnlock()
	endpointMap := endpointData.EndpointMap
	var gateways []*GateWayInfo
	for key, value := range endpointMap {
		name, ip, port := discovery.SplitRegisterKey(key)
		score := value.calculateScore()
		gateways = append(gateways, &GateWayInfo{
			Name:  name,
			Ip:    ip,
			Port:  port,
			Score: score,
		})
	}
	return gateways
}

func top5GateWays(gateways []*GateWayInfo) []*GateWayInfo {
	if len(gateways) > 5 {
		return gateways[:5]
	}
	return gateways
}
