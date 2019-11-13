// log包实现了基于logrus的日志管理器
package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"github.com/ebar-go/ego/library"
	"github.com/ebar-go/ego/http/constant"
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

type Context map[string]interface{}

const (
	DefaultAppChannel = "APP" // 日志的数据来源
	DefaultSystemChannel = "SYSTEM"
	DefaultRequestChannel = "REQUEST"
	DefaultMqChannel = "mq"
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

// NewContext 新的日志结构体
func NewContext(traceId string) Context {
	return Context{
		constant.TraceID : traceId,
	}
}


func (l *Logger) MergeSystemContext(context Context) Context {
	result := library.MergeMaps(Context{
		"service_name": l.systemParams.ServiceName,
		"service_port": l.systemParams.ServicePort,
	}, context)

	if _, ok := result[constant.TraceID]; ok {
		return result
	}

	result[constant.TraceID] = constant.TraceIdPrefix + library.UniqueId()
	return result
}


// Debug 调试等级
func (l *Logger) Debug(message string, context Context) {
	l.WithDefaultFields().WithField("context", l.MergeSystemContext(context)).Debug(message)
}

// Info 信息等级
func (l *Logger) Info(message string, context Context) {
	l.WithDefaultFields().WithField("context", l.MergeSystemContext(context)).Info(message)
}

// Warn 警告等级
func (l *Logger) Warn(message string, context Context) {
	l.WithDefaultFields().WithField("context", l.MergeSystemContext(context)).Warn(message)
}

// Error 错误等级
func (l *Logger) Error(message string, context Context) {
	l.WithDefaultFields().WithField("context", l.MergeSystemContext(context)).Error(message)
}

// Fatal 中断等级
func (l *Logger) Fatal(message string, context Context) {
	l.WithDefaultFields().WithField("context", l.MergeSystemContext(context)).Fatal(message)
}
