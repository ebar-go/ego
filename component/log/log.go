// log包实现了基于logrus的日志管理器
package log

import (
	"fmt"
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/ego/utils"
	"github.com/ebar-go/ego/utils/date"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// Logger 日志接口
type Logger interface {
	Info(message string, context Context)
	Debug(message string, context Context)
	Warn(message string, context Context)
	Error(message string, context Context)
	Fatal(message string, context Context)
	SetExtends(extends Context)
}

// logger 日志结构体
type logger struct {
	instance *logrus.Logger // logrus实例
	extends  Context
}

type Context map[string]interface{}

// New 获取默认的日志管理器，输出到控制台
func New() Logger {
	l := &logger{}
	l.instance = defaultInstance()
	return l
}

// NewFileLogger 根据文件初始化日志
func NewFileLogger(filePath string) Logger {
	logger := &logger{}
	logger.instance = defaultInstance()

	if !utils.PathExist(filePath) {
		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to init logger:%s,%s\n", filePath, err.Error())
			return logger
		}
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err == nil {
		logger.instance.Out = f
		fmt.Printf("Init Logger Success:%s\n", filePath)
	} else {
		fmt.Printf("Failed to init logger:%s,%s\n", filePath, err.Error())
	}

	return logger
}

// getDefaultLogInstance 实例化默认日志实例
func defaultInstance() *logrus.Logger {
	instance := logrus.New()

	// 设置日志格式为json
	instance.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "datetime",
			logrus.FieldKeyLevel: "level_name",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "caller",
		},
		TimestampFormat: date.TimeFormat,
	})
	instance.Level = logrus.DebugLevel

	return instance
}

func (l *logger) SetExtends(extends Context) {
	l.extends = extends
}

// withFields merge extends
func (l *logger) withFields(context Context) *logrus.Entry {
	if _, ok := context["trace_id"]; !ok {
		context["trace_id"] = trace.GetTraceId()
	}

	return l.instance.WithFields(logrus.Fields{
		"context": utils.MergeMaps(l.extends, context),
	})
}

// Debug
func (l *logger) Debug(message string, context Context) {
	l.withFields(context).Debug(message)
}

// Info
func (l *logger) Info(message string, context Context) {
	l.withFields(context).Info(message)
}

// Warn
func (l *logger) Warn(message string, context Context) {
	l.withFields(context).Warn(message)
}

// Error
func (l *logger) Error(message string, context Context) {
	l.withFields(context).Error(message)
}

// Fatal
func (l *logger) Fatal(message string, context Context) {
	l.withFields(context).Fatal(message)
}
