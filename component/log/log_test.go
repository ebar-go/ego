package log

import (
	"github.com/ebar-go/ego/utils/strings"
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
	traceId := strings.UUID()
	prepareLogger().Info("A group of walrus emerges from the ocean", Context{"trace_id": traceId, "hello": "world"})
}

// TestLogger_Info 测试Info
func TestLogger_Debug(t *testing.T) {
	traceId := strings.UUID()
	prepareLogger().Debug("A group of walrus emerges from the ocean", Context{"trace_id": traceId, "hello": "world"})
}

// TestLogger_Info 测试Info
func TestLogger_Warn(t *testing.T) {
	traceId := strings.UUID()
	prepareLogger().Warn("A group of walrus emerges from the ocean", Context{"trace_id": traceId, "hello": "world"})
}

// TestLogger_Info 测试Info
func TestLogger_Error(t *testing.T) {
	traceId := strings.UUID()
	prepareLogger().Error("A group of walrus emerges from the ocean", Context{"trace_id": traceId, "hello": "world"})
}

func TestNewManager(t *testing.T) {
	manager := NewManager(ManagerConf{
		SystemName: "test",
		SystemPort: 8080,
		LogPath:    "/tmp",
	})

	assert.NotNil(t, manager.App())
	assert.NotNil(t, manager.Mq())
	assert.NotNil(t, manager.System())
	assert.NotNil(t, manager.Request())
}
