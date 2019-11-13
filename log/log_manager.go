package log

import (
	"os"
	"fmt"
	"github.com/ebar-go/ego/library"
	"path/filepath"
	"path"
	"github.com/ebar-go/ego/http/constant"
	"github.com/sirupsen/logrus"
)

var systemLogManagerInstance *SystemLogManager

// SystemLogManager 系统管理器
type SystemLogManager struct {
	SystemName string
	SystemPort int
	AppDebug bool
	LogPath string
	app *Logger
	system *Logger
	request *Logger
	mq *Logger
	mqRotateDate string

}

func SetSystemLogManagerInstance(instance *SystemLogManager)  {
	systemLogManagerInstance = instance
}

// Init 初始化
func (manager *SystemLogManager) Init() {
	appPath := path.Join(manager.LogPath, manager.SystemName, constant.AppLogPrefix + manager.SystemName + constant.LogSuffix)
	systemPath := path.Join(manager.LogPath, manager.SystemName, constant.SystemLogPrefix + manager.SystemName + constant.LogSuffix)
	requestPath := path.Join(manager.LogPath, manager.SystemName, constant.RequestLogPrefix + manager.SystemName + constant.LogSuffix)

	manager.app = NewFileLogger(appPath)
	if !manager.AppDebug {
		manager.app.SetLevel(logrus.DebugLevel)
	}

	manager.app.SetSystemParam(LogSystemParam{
		ServiceName: manager.SystemName,
		ServicePort: manager.SystemPort,
		Channel: DefaultAppChannel,
	})

	manager.system = NewFileLogger(systemPath)
	manager.system.SetSystemParam(LogSystemParam{
		ServiceName: manager.SystemName,
		ServicePort: manager.SystemPort,
		Channel: DefaultSystemChannel,
	})

	manager.request = NewFileLogger(requestPath)
	manager.request.SetSystemParam(LogSystemParam{
		ServiceName: manager.SystemName,
		ServicePort: manager.SystemPort,
		Channel: DefaultRequestChannel,
	})

}

// App 应用日志
func App() *Logger {
	if systemLogManagerInstance.app == nil {
		systemLogManagerInstance.app = New()
	}

	return systemLogManagerInstance.app
}

// System 系统日志
func System() *Logger {
	if systemLogManagerInstance.app == nil {
		systemLogManagerInstance.app = New()
	}

	return systemLogManagerInstance.app
}

// Request 请求日志
func Request() *Logger {
	if systemLogManagerInstance.app == nil {
		systemLogManagerInstance.app = New()
	}

	return systemLogManagerInstance.app
}

// Mq 消息队列日志
func Mq() *Logger {

	currentDateStr := library.GetDateStr()
	if systemLogManagerInstance.mqRotateDate != currentDateStr {
		systemLogManagerInstance.mqRotateDate = currentDateStr
		mqPath := path.Join(systemLogManagerInstance.LogPath,
			DefaultMqChannel,
			systemLogManagerInstance.SystemName,
			systemLogManagerInstance.mqRotateDate + constant.LogSuffix)

		systemLogManagerInstance.mq = NewFileLogger(mqPath)
		systemLogManagerInstance.mq.SetSystemParam(LogSystemParam{
			ServiceName: systemLogManagerInstance.SystemName,
			ServicePort: systemLogManagerInstance.SystemPort,
			Channel: DefaultRequestChannel,
		})
	}

	if systemLogManagerInstance.mq == nil {
		systemLogManagerInstance.mq = New()
	}

	return systemLogManagerInstance.mq
}


// NewFileLogger 新的文件日志管理器
func NewFileLogger(filePath string) *Logger {
	er := New()

	if !library.IsPathExist(filePath) {
		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil{
			library.Debug(err)
			return er
		}
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err == nil {
		er.SetOutWriter(file)
		fmt.Printf("初始化日志文件:%s,成功\n", filePath)
	}else {
		fmt.Println("err:" + err.Error())
	}

	return er
}
