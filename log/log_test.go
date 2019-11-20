package log

import (
	"testing"
	"os"
	"fmt"
	"github.com/ebar-go/ego/helper"
	"github.com/stretchr/testify/assert"
)


// TestLogger_Init 测试初始化
func TestLogger_New(t *testing.T) {
	logger := New()
	assert.NotNil(t, logger)
}



func TestLogger_SetOutWriter(t *testing.T) {
	logger := New()
	filePath := helper.GetCurrentPath() + "/system.log"
	fmt.Println(filePath)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutWriter(file)
	}

	assert.Equal(t, file, logger.instance.Out)
}

// TestLogger_Info 测试Info
func TestLogger_Info(t *testing.T) {
	logger := New()
	logger.SetSystemParam(LogSystemParam{
		Channel: "REQUEST",
		ServiceName:"stock-manage",
		ServicePort:9523,
	})

	traceId := helper.UniqueId()
	logger.Info("A group of walrus emerges from the ocean", Context{"trace_id": traceId})
}


func TestLoggerFile(t *testing.T)  {
	var err error
	logger := New()


	filePath := helper.GetCurrentPath() + "/system.log"
	fmt.Println(filePath)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutWriter(file)
	}else {
		fmt.Println("err:" + err.Error())
	}

}

