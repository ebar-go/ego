package log

import (
	"github.com/ebar-go/ego/helper"
	"path"
	"sync"
)

const (
	DefaultLogPath = "/wwwlogs/"
	SystemLogPrefix = "system_"
	RequestLogPrefix = "request_"
	AppLogPrefix = "app_"
	MqLogPrefix = "mq_"
	LogSuffix = ".log"

	AppLogComponentName = "app"
	TraceLogComponentName = "trace"
	MqLogComponentName = "mq"
	SystemLogComponentName = "phplogs"
)

var m Manager

// Manager 接口
type Manager interface {
	rotate() *manager
}

// InitManager 初始化管理器
func InitManager(conf ManagerConf)  {
	conf.LogPath = helper.DefaultString(conf.LogPath, DefaultLogPath)
	instance := &manager{conf:conf}
	instance.rotate()

	m = instance
}

// manager 系统日志管理器
type manager struct {
	sync.Mutex
	conf ManagerConf
	rotateDate string // 分割日期
	app Logger
	system Logger
	request Logger
	mq Logger
}

// InitAppLogger
func (manager *manager) InitAppLogger()  {
	filePath := manager.GetLogPath(AppLogComponentName, AppLogPrefix)
	manager.app = NewFileLogger(filePath)
}

// InitRequestLogger
func (manager *manager) InitRequestLogger()  {
	filePath := manager.GetLogPath(TraceLogComponentName, RequestLogPrefix)
	manager.request = NewFileLogger(filePath)
}

// InitSystemLogger
func (manager *manager) InitSystemLogger()  {
	filePath := manager.GetLogPath(SystemLogComponentName, SystemLogPrefix)
	manager.system = NewFileLogger(filePath)
}

// GetLogPath 根据组件名称获取日志路径
func (manager *manager) GetLogPath(componentName, prefix string) string {
	return path.Join(manager.conf.LogPath,
		componentName,
		manager.conf.SystemName,
		prefix + manager.rotateDate + LogSuffix)
}

// InitMqLogger
func (manager *manager) InitMqLogger()  {
	filePath := manager.GetLogPath(MqLogComponentName, MqLogPrefix)
	manager.mq = NewFileLogger(filePath)
}

// ManagerConf 日志配置
type ManagerConf struct {
	SystemName string
	SystemPort int
	LogPath string
}

// rotate 分割日志文件
func (m *manager) rotate() *manager {
	currentDateStr := helper.GetDateStr()
	if currentDateStr != m.rotateDate {
		m.Lock()
		defer m.Unlock()

		m.rotateDate = currentDateStr

		m.InitAppLogger()
		m.InitRequestLogger()
		m.InitSystemLogger()
		m.InitMqLogger()
	}

	return m
}

// App
func App() Logger  {
	return m.rotate().app
}

// Request
func Request() Logger  {
	return m.rotate().request
}

// System
func System() Logger  {
	return m.rotate().system
}

// Mq
func Mq() Logger  {
	return m.rotate().mq
}
