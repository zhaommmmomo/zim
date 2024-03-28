package log

import (
	"context"
	"github.com/zhaommmmomo/zim/common/trace"
	"go.uber.org/zap"
)

func CtxTid(ctx *context.Context) zap.Field {
	return zap.String(trace.TID, trace.GetTraceId(ctx))
}

func Tid(tid string) zap.Field {
	return zap.String(trace.TID, tid)
}

func Err(err error) zap.Field {
	return zap.Error(err)
}

func String(key, value string) zap.Field {
	return zap.String(key, value)
}

func Any(key string, value any) zap.Field {
	return zap.Any(key, value)
}
