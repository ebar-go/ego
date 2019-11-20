package log

import (
	"github.com/ebar-go/ego/library"
	"path"
	"github.com/ebar-go/ego/http/constant"
	"sync"
)

var manager = new(Manager)
var once = new(sync.Once)

// Manager 系统日志管理器
type Manager struct {
	conf ManagerConf
	rotateDate string // 分割日期
	app *Logger
	system *Logger
	request *Logger
	mq *Logger
}

// GetSystemParam
func (conf ManagerConf) GetSystemParam() SystemParam {
	return SystemParam{
		ServiceName: conf.SystemName,
		ServicePort: conf.SystemPort,
	}
}

// InitAppLogger
func (manager *Manager) InitAppLogger()  {
	filePath := path.Join(manager.conf.LogPath,
		constant.AppLogComponentName,
		manager.conf.SystemName,
		constant.AppLogPrefix + manager.rotateDate + constant.LogSuffix)
	manager.app = NewFileLogger(filePath)
	manager.app.SetSystemParam(manager.conf.GetSystemParam())
}

// InitRequestLogger
func (manager *Manager) InitRequestLogger()  {
	filePath := path.Join(manager.conf.LogPath,
		constant.TraceLogComponentName,
		manager.conf.SystemName,
		constant.RequestLogPrefix + manager.rotateDate + constant.LogSuffix)
	manager.request = NewFileLogger(filePath)
	manager.request.SetSystemParam(manager.conf.GetSystemParam())
}

// InitSystemLogger
func (manager *Manager) InitSystemLogger()  {
	filePath := path.Join(manager.conf.LogPath,
		constant.SystemLogComponentName,
		manager.conf.SystemName,
		constant.SystemLogPrefix + manager.rotateDate + constant.LogSuffix)
	manager.system = NewFileLogger(filePath)
	manager.system.SetSystemParam(manager.conf.GetSystemParam())
}

// InitMqLogger
func (manager *Manager) InitMqLogger()  {
	filePath := path.Join(manager.conf.LogPath,
		constant.MqLogComponentName,
		manager.conf.SystemName,
		constant.MqLogPrefix + manager.rotateDate + constant.LogSuffix)
	manager.mq = NewFileLogger(filePath)
	manager.mq.SetSystemParam(manager.conf.GetSystemParam())
}

// ManagerConf 日志配置
type ManagerConf struct {
	SystemName string
	SystemPort int
	AppDebug bool
	LogPath string
}

// InitManager 初始化日志管理器
func InitManager(conf ManagerConf) {
	once.Do(func() {
		manager.conf = conf
		manager.rotate()
	})

}

// rotate 分割日志文件
func (m *Manager) rotate() *Manager {
	currentDateStr := library.GetDateStr()
	if currentDateStr != m.rotateDate {
		// TODO 需要验证并发时的安全性
		m.rotateDate = currentDateStr

		m.InitAppLogger()
		m.InitRequestLogger()
		m.InitSystemLogger()
		m.InitMqLogger()
	}

	return m
}

// App 应用日志
func App() *Logger {
	return manager.rotate().app
}

// System 系统日志
func System() *Logger {
	return manager.rotate().system
}

// Request 请求日志
func Request() *Logger {
	return manager.rotate().request
}

// Mq 消息队列日志
func Mq() *Logger {
	return manager.rotate().mq
}

