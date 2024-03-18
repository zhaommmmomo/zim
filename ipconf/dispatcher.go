package ipconf

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zhaommmmomo/zim/common/domain"
)

func GetEndpoints(context context.Context, appCtx *app.RequestContext) {
	defer func() {
		if err := recover(); err != nil {
			appCtx.JSON(consts.StatusBadRequest, utils.H{"err": err})
		}
	}()
	// 构建请求上下文参数
	ctx := domain.BuildIpConfContext(&context, appCtx)
	// 获取集群中对应的endpoints  按照对应的分数进行排序了
	endpoints := dispatch(ctx)
	// 返回数据
	appCtx.JSON(consts.StatusOK, domain.PackRes(top5Endpoints(endpoints)))
}

func dispatch(ctx *domain.Context) []*domain.Endpoint {
	return []*domain.Endpoint{{
		Ip:    "127.0.0.1",
		Port:  9000,
		Score: 11,
	}, {
		Ip:    "127.0.0.1",
		Port:  9001,
		Score: 9,
	}, {
		Ip:    "127.0.0.1",
		Port:  9002,
		Score: 5,
	}, {
		Ip:    "127.0.0.1",
		Port:  9003,
		Score: 1,
	}, {
		Ip:    "127.0.0.1",
		Port:  9004,
		Score: 20,
	}, {
		Ip:    "127.0.0.1",
		Port:  9005,
		Score: 31,
	}}
}

func top5Endpoints(endpoints []*domain.Endpoint) []*domain.Endpoint {
	if len(endpoints) > 5 {
		return endpoints[:5]
	}
	return endpoints
}
