package trace

import (
	"context"
	"github.com/zhaommmmomo/zim/common/utils"
)

const TID = "TID"

func NewCustomCtxWithTraceId(prefix string) *context.Context {
	ctx := context.WithValue(context.Background(), TID, prefix+"-"+utils.NewUUID())
	return &ctx
}

func NewCtxWithTraceId() *context.Context {
	ctx := context.WithValue(context.Background(), TID, utils.NewUUID())
	return &ctx
}

func GetTraceId(ctx *context.Context) string {
	return utils.ConvertToString((*ctx).Value(TID))
}
