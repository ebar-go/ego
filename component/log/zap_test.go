package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestZap(t *testing.T) {
	logger := NewZap("/tmp/zap/request.log", zap.DebugLevel, zap.String("system_name", "app"))

	defer logger.Sync()

	logger.Info("info", zap.Field{
		Key:       "hello",
		String:    "world",
		Type: zapcore.StringType,
	})
	logger.Error("error", zap.Field{
		Key:       "hello",
		String:    "world",
		Type: zapcore.StringType,
	})
	logger.DPanic("panic", zap.Field{
		Key:       "hello",
		String:    "world",
		Type: zapcore.StringType,
	})
}
