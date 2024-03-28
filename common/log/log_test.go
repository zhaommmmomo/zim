package log

import (
	"github.com/zhaommmmomo/zim/common/domain"
	"github.com/zhaommmmomo/zim/common/trace"
	"testing"
)

func TestLog(t *testing.T) {
	ctx := trace.NewCtxWithTraceId()
	InfoCtx(ctx, "info")
	WarnCtx(ctx, "warn")
	ErrorCtx(ctx, "err")

	message := &domain.Message{
		FHeader: &domain.FixedHeader{
			V:          1,
			Cmd:        2,
			VarHLen:    15,
			PayloadLen: 12,
			Crc32sum:   0,
		},
		VHeader: []byte("Variable Header"),
		Payload: []byte("Payload Data"),
	}

	InfoCtx(ctx, "test", Any("data", message))
}
