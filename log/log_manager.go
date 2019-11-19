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
	rotateDate string
}

func SetSystemLogManagerInstance(instance *SystemLogManager)  {
	systemLogManagerInstance = instance
}

func (manager *SystemLogManager) validate() *SystemLogManager {
	currentDateStr := library.GetDateStr()
	if currentDateStr != manager.rotateDate {
		manager.rotateDate = currentDateStr
		manager.initLogByRotate()
	}

	return manager
}

func (manager *SystemLogManager) initLogByRotate()  {
	systemParam := SystemParam{
		ServiceName: manager.SystemName,
		ServicePort: manager.SystemPort,
	}

	appPath := path.Join(manager.LogPath,
		constant.AppLogComponentName,
		manager.SystemName,
		constant.AppLogPrefix + manager.rotateDate + constant.LogSuffix)
	manager.app = NewFileLogger(appPath)
	manager.app.SetSystemParam(systemParam)

	requestPath := path.Join(manager.LogPath,
		constant.TraceLogComponentName,
		manager.SystemName,
		constant.RequestLogPrefix + manager.rotateDate + constant.LogSuffix)
	manager.request = NewFileLogger(requestPath)
	manager.request.SetSystemParam(systemParam)

	systemPath := path.Join(manager.LogPath,
		constant.SystemLogComponentName,
		manager.SystemName,
		constant.SystemLogPrefix + manager.rotateDate + constant.LogSuffix)
	manager.system = NewFileLogger(systemPath)
	manager.system.SetSystemParam(systemParam)

	mqPath := path.Join(manager.LogPath,
		constant.MqLogComponentName,
		manager.SystemName,
		constant.MqLogPrefix + manager.rotateDate + constant.LogSuffix)
	manager.mq = NewFileLogger(mqPath)
	manager.mq.SetSystemParam(systemParam)
}


// Init 初始化
func (manager *SystemLogManager) Init() {
	manager.rotateDate = library.GetDateStr()
	manager.initLogByRotate()
}

// App 应用日志
func App() *Logger {
	if systemLogManagerInstance.app == nil {
		systemLogManagerInstance.app = New()
	}
	return systemLogManagerInstance.validate().app
}

// System 系统日志
func System() *Logger {
	if systemLogManagerInstance.system == nil {
		systemLogManagerInstance.system = New()
	}
	return systemLogManagerInstance.validate().system
}

// Request 请求日志
func Request() *Logger {
	if systemLogManagerInstance.request == nil {
		systemLogManagerInstance.request = New()
	}
	return systemLogManagerInstance.validate().request
}

// Mq 消息队列日志
func Mq() *Logger {
	if systemLogManagerInstance.mq == nil {
		systemLogManagerInstance.mq = New()
	}

	return systemLogManagerInstance.validate().mq
}

