// log包实现了基于logrus的日志管理器
package log

import (
	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	"io"
	"os"
)

// Logger 日志结构体
//   - Key 支持自定义日志的字段名称，默认是LoggerDefaultKey
//   - Out 日志的输出目标,可以是文件,也可以是os.Stdout
type Logger struct {
	initialize bool // 仅初始化一次
	instance   *logrus.Logger // logrus实例
	Key string
	Out io.Writer 
}

var systemLogger *Logger

const (
	LoggerDefaultKey = "title" // 日志的默认字段名称 title : this is content
)

// DefaultLogger 默认的日志管理器，输出到控制台
func DefaultLogger() *Logger {
	logger := &Logger{
		Out: os.Stdout,
	}

	logger.Init()
	return logger
}

// SetSystemLogger 设置系统日志
func SetSystemLogger(logger *Logger)  {
	systemLogger = logger
}

// GetSystemLogger 获取系统日志
func GetSystemLogger() *Logger {
	if systemLogger == nil {
		systemLogger = DefaultLogger()
	}

	return systemLogger
}

// Init 初始化日志管理器
func (logger *Logger) Init() error {
	// 拒绝重复初始化
	if logger.initialize {
		return errors.New("请勿重复初始化")
	}

	logger.instance = logrus.New()

	// 设置日志格式为json
	logger.instance.SetFormatter(&logrus.JSONFormatter{})
	logger.instance.Out = logger.Out

	//设置最低log level
	//logger.instance.SetLevel(logrus.InfoLevel)

	// 如果没有传key的键名，则使用默认值
	if logger.Key == "" {
		logger.Key = LoggerDefaultKey
	}

	logger.initialize = true
	return nil
}

// Debug 调试等级,记录title为日志备注，context为日志的message内容
func (l *Logger) Debug(title string, context ...interface{}) {
	l.instance.WithField(l.Key, title).Debug(context...)
}

// Info 信息等级,记录title为日志备注，context为日志的message内容
func (l *Logger) Info(title string, context ...interface{}) {
	l.instance.WithField(l.Key, title).Info(context...)
}

// Warn 警告等级,记录title为日志备注，context为日志的message内容
func (l *Logger) Warn(title string, context ...interface{}) {
	l.instance.WithField(l.Key, title).Warn(context...)
}

// Error 错误等级,记录title为日志备注，context为日志的message内容
func (l *Logger) Error(title string, context ...interface{}) {
	l.instance.WithField(l.Key, title).Error(context...)
}

// Fatal 中断等级,记录title为日志备注，context为日志的message内容
func (l *Logger) Fatal(title string, context ...interface{}) {
	l.instance.WithField(l.Key, title).Fatal(context...)
}
