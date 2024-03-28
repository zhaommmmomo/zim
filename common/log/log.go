package log

import (
	"context"
	"github.com/natefinch/lumberjack"
	"github.com/zhaommmmomo/zim/common/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	logger *zap.Logger
)

func Init() {
	var core zapcore.Core
	if config.IsDebug() {
		core = zapcore.NewCore(
			initEncoder(),
			initWriteSyncer(),
			zap.DebugLevel,
		)
	} else {
		core = zapcore.NewCore(
			initEncoder(),
			initWriteSyncer(),
			zap.InfoLevel,
		)
	}
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel))
	return
}

func initEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	return encoder
}

func initWriteSyncer() zapcore.WriteSyncer {
	ws := zapcore.NewMultiWriteSyncer(zapcore.AddSync(&lumberjack.Logger{
		Filename: "../../logs/zim.log",
	}), zapcore.AddSync(os.Stderr))
	return ws
}

func DebugCtx(ctx *context.Context, s string, fields ...zap.Field) {
	logger.Debug(s, append(fields, CtxTid(ctx))...)
}

func InfoCtx(ctx *context.Context, s string, fields ...zap.Field) {
	logger.Info(s, append(fields, CtxTid(ctx))...)
}

func WarnCtx(ctx *context.Context, s string, fields ...zap.Field) {
	logger.Warn(s, append(fields, CtxTid(ctx))...)
}

func ErrorCtx(ctx *context.Context, s string, fields ...zap.Field) {
	logger.Error(s, append(fields, CtxTid(ctx))...)
}

func Debug(s string, fields ...zap.Field) {
	logger.Debug(s, fields...)
}

func Info(s string, fields ...zap.Field) {
	logger.Info(s, fields...)
}

func Warn(s string, fields ...zap.Field) {
	logger.Warn(s, fields...)
}

func Error(s string, fields ...zap.Field) {
	logger.Error(s, fields...)
}

func Sync() {
	logger.Sync()
}
