package log

import (
	"github.com/ebar-go/ego/library"
	"path"
	"github.com/ebar-go/ego/http/constant"
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
	manager.app.SetSystemParam(SystemParam{
		ServiceName: manager.SystemName,
		ServicePort: manager.SystemPort,
		Channel: DefaultAppChannel,
	})

	manager.system = NewFileLogger(systemPath)
	manager.system.SetSystemParam(SystemParam{
		ServiceName: manager.SystemName,
		ServicePort: manager.SystemPort,
		Channel: DefaultSystemChannel,
	})

	manager.request = NewFileLogger(requestPath)
	manager.request.SetSystemParam(SystemParam{
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
	if systemLogManagerInstance.request == nil {
		systemLogManagerInstance.request = New()
	}

	return systemLogManagerInstance.request
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
		systemLogManagerInstance.mq.SetSystemParam(SystemParam{
			ServiceName: systemLogManagerInstance.SystemName,
			ServicePort: systemLogManagerInstance.SystemPort,
			Channel: DefaultMqChannel,
		})
	}

	if systemLogManagerInstance.mq == nil {
		systemLogManagerInstance.mq = New()
	}

	return systemLogManagerInstance.mq
}

