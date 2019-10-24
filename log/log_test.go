package log

import (
	"testing"
	"os"
	"github.com/ebar-go/ego/test"
)

// 准备日志管理器
func getLogger() *Logger {
	logger := &Logger{
		Out: os.Stdout,
	}

	logger.Init()
	return logger
}

// TestLogger_Init 测试初始化
func TestLogger_Init(t *testing.T) {
	logger := &Logger{
		Out: os.Stdout,
	}
	err := logger.Init()
	test.AssertNil(t, err)
}

// TestLogger_Info 测试Info
func TestLogger_Info(t *testing.T) {
	getLogger().Info("test")
}

// TestLogger_Debug 测试Debug
func TestLogger_Debug(t *testing.T) {
	getLogger().Debug("test debug", 123, 456)
}

// TestLogger_Warn 测试Warn
func TestLogger_Warn(t *testing.T) {
	getLogger().Warn("test warn", 123, 456)
}

// TestLogger_Error 测试Error
func TestLogger_Error(t *testing.T) {
	getLogger().Error("test error", 123, 456)
}

// TestLogger_Fatal 测试Fatal
func TestLogger_Fatal(t *testing.T) {
	t.SkipNow()
	getLogger().Fatal("test fatal", 123, 456)
}


