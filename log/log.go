// log包实现了基于logrus的日志管理器
package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"github.com/ebar-go/ego/library"
)

// Logger 日志结构体
type Logger struct {
	instance   *logrus.Logger // logrus实例
	systemParams LogSystemParam
}

type LogSystemParam struct {
	ServiceName string `json:"service_name"`
	ServicePort int `json:"service_port"`
	Channel string `json:"channel"`
}

type LogContextInterface map[string]interface{}

const (
	DefaultAppChannel = "APP" // 日志的数据来源
	DefaultSystemChannel = "SYSTEM"
	DefaultRequestChannel = "REQUEST"
)

// New 获取默认的日志管理器，输出到控制台
func New() *Logger {
	l := &Logger{}
	l.instance = getDefaultLogInstance()
	return l
}

// GetInstance 获取logrus实例
func (l *Logger) GetInstance() *logrus.Logger  {
	return l.instance
}

// NewContext
func (l *Logger) NewContext(traceId string ) LogContextInterface {
	context := LogContextInterface{
		"service_name": l.systemParams.ServiceName,
		"service_port": l.systemParams.ServicePort,
		"trace_id" : traceId,
	}

	return context
}


// SetLevel 设置日志等级
func (l *Logger) SetLevel(level logrus.Level) {
	l.instance.Level = level
}

// SetOutWriter 设置输出,可以是文件，也可以是os.StdOut
func (l *Logger) SetOutWriter(out io.Writer)  {
	l.instance.Out = out
}

// SetSystemParam 设置系统参数
func (l *Logger) SetSystemParam(param LogSystemParam)  {
	l.systemParams = param
}


// getDefaultLogInstance 实例化默认日志实例
func getDefaultLogInstance() *logrus.Logger {
	instance := logrus.New()

	// 设置日志格式为json
	instance.SetFormatter(&logrus.JSONFormatter{
		  	FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "datetime",
				logrus.FieldKeyLevel: "level_name",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
		   },
		  	TimestampFormat: library.GetDefaultTimeFormat(),
	})


	return instance
}

func (l *Logger) WithDefaultFields() *logrus.Entry {
	if l.systemParams.Channel == "" {
		l.systemParams.Channel = DefaultAppChannel
	}

	return l.instance.WithFields(logrus.Fields{
		"channel": l.systemParams.Channel,
	})
}


// Debug 调试等级
func (l *Logger) Debug(message string, context LogContextInterface) {
	l.WithDefaultFields().WithField("context", context).Debug(context)
}

// Info 信息等级
func (l *Logger) Info(message string, context LogContextInterface) {
	l.WithDefaultFields().WithField("context", context).Info(message)
}

// Warn 警告等级
func (l *Logger) Warn(message string, context LogContextInterface) {
	l.WithDefaultFields().Warn(context)
}

// Error 错误等级
func (l *Logger) Error(message string, context LogContextInterface) {
	l.WithDefaultFields().Error(context)
}

// Fatal 中断等级
func (l *Logger) Fatal(message string, context LogContextInterface) {
	l.WithDefaultFields().Fatal(context)
}
