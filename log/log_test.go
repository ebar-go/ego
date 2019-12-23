package log

import (
	"github.com/ebar-go/ego/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNew 测试初始化
func TestNew(t *testing.T) {
	logger := New()
	assert.NotNil(t, logger)
}

func TestNewFileLogger(t *testing.T) {
	filePath := "/tmp/system.log"
	logger := NewFileLogger(filePath)
	assert.NotNil(t, logger)
}

func prepareLogger() Logger {
	return NewFileLogger("/tmp/test.log")
}

// TestLogger_Info 测试Info
func TestLogger_Info(t *testing.T) {
	traceId := helper.UniqueId()
	prepareLogger().Info("A group of walrus emerges from the ocean", Context{"trace_id": traceId, "hello": "world"})
}

// TestLogger_Info 测试Info
func TestLogger_Debug(t *testing.T) {
	traceId := helper.UniqueId()
	prepareLogger().Debug("A group of walrus emerges from the ocean", Context{"trace_id": traceId, "hello": "world"})
}

// TestLogger_Info 测试Info
func TestLogger_Warn(t *testing.T) {
	traceId := helper.UniqueId()
	prepareLogger().Warn("A group of walrus emerges from the ocean", Context{"trace_id": traceId, "hello": "world"})
}

// TestLogger_Info 测试Info
func TestLogger_Error(t *testing.T) {
	traceId := helper.UniqueId()
	prepareLogger().Error("A group of walrus emerges from the ocean", Context{"trace_id": traceId, "hello": "world"})
}

func TestInitManager(t *testing.T) {
	InitManager(ManagerConf{
		SystemName: "test",
		SystemPort: 8080,
		LogPath:    "/tmp",
	})

	assert.NotNil(t, App())
	assert.NotNil(t, Mq())
	assert.NotNil(t, System())
	assert.NotNil(t, Request())
}
