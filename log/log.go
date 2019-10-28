// log包实现了基于logrus的日志管理器
package log

import (
	"github.com/sirupsen/logrus"
	"io"
)

// Logger 日志结构体
//   - key 支持自定义日志的字段名称，默认是LoggerDefaultkey
//   - Out 日志的输出目标,可以是文件,也可以是os.Stdout
type Logger struct {
	instance   *logrus.Logger // logrus实例
	key string
}

var systemLogger *Logger

const (
	defaultKey = "title" // 日志的默认字段名称 title : this is content
)

// New 获取默认的日志管理器，输出到控制台
func New() *Logger {
	l := &Logger{}
	l.instance = getDefaultLogInstance()
	l.key = defaultKey

	return l
}

// SetSystemLogger 设置系统日志
func SetSystemLogger(logger *Logger)  {
	systemLogger = logger
}

// GetSystemLogger 获取系统日志
func GetSystemLogger() *Logger {
	if systemLogger == nil {
		systemLogger = New()
	}

	return systemLogger
}

// SetOutWriter 设置输出,可以是文件，也可以是os.StdOut
func (l *Logger) SetOutWriter(out io.Writer)  {
	l.instance.Out = out
}

// SetKey 设置字段名称
func (l *Logger) SetKey(key string)  {
	l.key = key
}

// getDefaultLogInstance 实例化默认日志实例
func getDefaultLogInstance() *logrus.Logger {
	instance := logrus.New()

	// 设置日志格式为json
	instance.SetFormatter(&logrus.JSONFormatter{})
	return instance
}

// Debug 调试等级,记录title为日志备注，context为日志的message内容
func (l *Logger) Debug(title string, context ...interface{}) {
	l.instance.WithField(l.key, title).Debug(context...)
}

// Info 信息等级,记录title为日志备注，context为日志的message内容
func (l *Logger) Info(title string, context ...interface{}) {
	l.instance.WithField(l.key, title).Info(context...)
}

// Warn 警告等级,记录title为日志备注，context为日志的message内容
func (l *Logger) Warn(title string, context ...interface{}) {
	l.instance.WithField(l.key, title).Warn(context...)
}

// Error 错误等级,记录title为日志备注，context为日志的message内容
func (l *Logger) Error(title string, context ...interface{}) {
	l.instance.WithField(l.key, title).Error(context...)
}

// Fatal 中断等级,记录title为日志备注，context为日志的message内容
func (l *Logger) Fatal(title string, context ...interface{}) {
	l.instance.WithField(l.key, title).Fatal(context...)
}
