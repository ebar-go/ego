package log

import (
	"fmt"
	"github.com/ebar-go/ego/utils/date"
	"github.com/ebar-go/ego/utils/strings"
	"path"
	"sync"
)

const (
	DefaultLogPath   = "/wwwlogs/"
	SystemLogPrefix  = "system_"
	RequestLogPrefix = "request_"
	AppLogPrefix     = "app_"
	MqLogPrefix      = "mq_"
	LogSuffix        = ".log"

	AppLogComponentName    = "app"
	TraceLogComponentName  = "trace"
	MqLogComponentName     = "mq"
	SystemLogComponentName = "phplogs"
)

// Manager 接口
type Manager interface {
	App() Logger
	Request() Logger
	System() Logger
	Mq() Logger
	rotate() *manager
}

// InitManager 初始化管理器
func NewManager(conf ManagerConf) Manager {
	conf.LogPath = strings.Default(conf.LogPath, DefaultLogPath)
	instance := &manager{conf: conf}
	instance.rotate()

	return instance
}

// manager 系统日志管理器
type manager struct {
	sync.Mutex
	conf       ManagerConf
	rotateDate string // 分割日期
	app        Logger
	system     Logger
	request    Logger
	mq         Logger
}

// InitAppLogger
func (manager *manager) InitAppLogger() {
	manager.app = manager.getLogger(AppLogComponentName, AppLogPrefix)
}

// InitRequestLogger
func (manager *manager) InitRequestLogger() {
	manager.request = manager.getLogger(TraceLogComponentName, RequestLogPrefix)
}

// InitSystemLogger
func (manager *manager) InitSystemLogger() {
	manager.system = manager.getLogger(SystemLogComponentName, SystemLogPrefix)
}

// GetLogPath 根据组件名称获取日志路径
func (manager *manager) getLogger(componentName, prefix string) Logger {
	filePath := path.Join(manager.conf.LogPath,
		componentName,
		manager.conf.SystemName,
		prefix+manager.rotateDate+LogSuffix)

	fmt.Println(componentName, filePath)
	logger := NewFileLogger(filePath)
	logger.SetExtends(Context{
		"service_name": manager.conf.SystemName,
		"service_port": manager.conf.SystemPort,
	})
	return logger
}

// InitMqLogger
func (manager *manager) InitMqLogger() {
	manager.mq = manager.getLogger(MqLogComponentName, MqLogPrefix)
}

// ManagerConf 日志配置
type ManagerConf struct {
	SystemName string
	SystemPort int
	LogPath    string
}

// rotate 分割日志文件
func (manager *manager) rotate() *manager {
	currentDateStr := date.GetDateStr()
	if currentDateStr != manager.rotateDate {
		manager.Lock()
		defer manager.Unlock()

		manager.rotateDate = currentDateStr

		manager.InitAppLogger()
		manager.InitRequestLogger()
		manager.InitSystemLogger()
		manager.InitMqLogger()
	}

	return manager
}

// App
func (m *manager) App() Logger {
	return m.rotate().app
}

// Request
func (m *manager) Request() Logger {
	return m.rotate().request
}

// System
func (m *manager) System() Logger {
	return m.rotate().system
}

// Mq
func (m *manager) Mq() Logger {
	return m.rotate().mq
}
