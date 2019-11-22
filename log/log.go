// log包实现了基于logrus的日志管理器
package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"github.com/ebar-go/ego/helper"
	"os"
	"path/filepath"
	"fmt"
	"github.com/ebar-go/ego/component/trace"
)

// Logger 日志结构体
type Logger struct {
	instance     *logrus.Logger // logrus实例
	systemParams SystemParam
}

type SystemParam struct {
	ServiceName string `json:"service_name"`
	ServicePort int    `json:"service_port"`
}

type Context map[string]interface{}

// New 获取默认的日志管理器，输出到控制台
func New() *Logger {
	l := &Logger{}
	l.instance = getDefaultLogInstance()
	return l
}


// NewFileLogger 新的文件日志管理器
func NewFileLogger(filePath string) *Logger {
	logger := New()

	if !helper.IsPathExist(filePath) {
		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil{
			helper.Debug(err)
			return logger
		}
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err == nil {
		logger.SetOutWriter(file)
		fmt.Printf("Init Logger Success:%s\n", filePath)
	}else {
		fmt.Printf("Failed to init logger:%s,%s\n", filePath ,err.Error())
	}

	return logger
}


// SetLevel 设置日志等级
func (l *Logger) SetLevel(level logrus.Level) {
	l.instance.Level = level
}

// SetOutWriter 设置输出,可以是文件，也可以是os.StdOut
func (l *Logger) SetOutWriter(out io.Writer) {
	l.instance.Out = out
}

// SetSystemParam 设置系统参数
func (l *Logger) SetSystemParam(param SystemParam) {
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
		TimestampFormat: helper.GetDefaultTimeFormat(),
	})
	instance.Level = logrus.DebugLevel

	return instance
}

// withFields 携带字段
func (l *Logger) withFields(context Context) *logrus.Entry {

	if _, ok := context["trace_id"]; !ok {
		context["trace_id"] = trace.GetTraceId()
	}

	return l.instance.WithFields(logrus.Fields{
		"context": helper.MergeMaps(Context{
			"service_name": l.systemParams.ServiceName,
			"service_port": l.systemParams.ServicePort,
		}, context),
	})
}

// Debug 调试等级
func (l *Logger) Debug(message string, context Context) {
	l.withFields(context).Debug(message)
}

// Info 信息等级
func (l *Logger) Info(message string, context Context) {
	l.withFields(context).Info(message)
}

// Warn 警告等级
func (l *Logger) Warn(message string, context Context) {
	l.withFields(context).Warn(message)
}

// Error 错误等级
func (l *Logger) Error(message string, context Context) {
	l.withFields(context).Error(message)
}

// Fatal 中断等级
func (l *Logger) Fatal(message string, context Context) {
	l.withFields(context).Fatal(message)
}
