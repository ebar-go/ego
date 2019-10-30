package log

import (
	"testing"
	"os"
	"fmt"
	"github.com/ebar-go/ego/library"
	"github.com/stretchr/testify/assert"
)


// TestLogger_Init 测试初始化
func TestLogger_New(t *testing.T) {
	logger := New()
	assert.NotNil(t, logger)
}

func TestGetSystemLogger(t *testing.T) {
	assert.NotNil(t, GetSystemLogger())
}

func TestSetSystemLogger(t *testing.T) {
	SetSystemLogger(New())
}

func TestLogger_SetKey(t *testing.T) {
	logger := New()
	key := "name"
	logger.SetKey(key)
	assert.Equal(t, key, logger.key)
}

func TestLogger_SetOutWriter(t *testing.T) {
	logger := New()
	filePath := library.GetCurrentPath() + "/system.log"
	fmt.Println(filePath)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutWriter(file)
	}

	assert.Equal(t, file, logger.instance.Out)
}

// TestLogger_Info 测试Info
func TestLogger_Info(t *testing.T) {
	New().Info("test")
}

// TestLogger_Debug 测试Debug
func TestLogger_Debug(t *testing.T) {
	New().Debug("test debug", 123, 456)
}

// TestLogger_Warn 测试Warn
func TestLogger_Warn(t *testing.T) {
	New().Warn("test warn", 123, 456)
}

// TestLogger_Error 测试Error
func TestLogger_Error(t *testing.T) {
	New().Error("test error", 123, 456)
}

// TestLogger_Fatal 测试Fatal
func TestLogger_Fatal(t *testing.T) {
	t.SkipNow()
	New().Fatal("test fatal", 123, 456)
}

func TestLoggerFile(t *testing.T)  {
	var err error
	logger := New()


	filePath := library.GetCurrentPath() + "/system.log"
	fmt.Println(filePath)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutWriter(file)
	}else {
		fmt.Println("err:" + err.Error())
	}

	logger.Info("test info", 123, 456)
}

