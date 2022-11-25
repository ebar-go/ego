package logger

import (
	"fmt"
	"log"
	"sync"
)

var (
	Info   = Default().Info
	Infof  = Default().Infof
	Error  = Default().Error
	Errorf = Default().Errorf
	Debug  = Default().Debug
	Debugf = Default().Debugf
	Warn   = Default().Warn
	Warnf  = Default().Warnf
	Fatal  = Default().Fatal
	Fatalf = Default().Fatalf
)

type Printer interface {
	Println(v ...any)
}

type Logger struct {
	instance Printer
}

type LogLevel string

const (
	LevelInfo  LogLevel = "info"
	LevelDebug LogLevel = "debug"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
	LevelFatal LogLevel = "fatal"
	LevelPanic LogLevel = "panic"
)

func (l *Logger) Log(level LogLevel, message string) { l.instance.Println(level, message) }
func (l *Logger) Logf(level LogLevel, format string, a ...interface{}) {
	l.instance.Println(level, fmt.Sprintf(format, a...))
}

func (l *Logger) Info(message string)  { l.Log(LevelInfo, message) }
func (l *Logger) Debug(message string) { l.Log(LevelDebug, message) }
func (l *Logger) Warn(message string)  { l.Log(LevelWarn, message) }
func (l *Logger) Error(message string) { l.Log(LevelError, message) }
func (l *Logger) Fatal(message string) { l.Log(LevelFatal, message) }
func (l *Logger) Panic(message string) { l.Log(LevelPanic, message) }

func (l *Logger) Infof(format string, args ...interface{})  { l.Logf(LevelInfo, format, args...) }
func (l *Logger) Debugf(format string, args ...interface{}) { l.Logf(LevelDebug, format, args...) }
func (l *Logger) Warnf(format string, args ...interface{})  { l.Logf(LevelWarn, format, args...) }
func (l *Logger) Errorf(format string, args ...interface{}) { l.Logf(LevelError, format, args...) }
func (l *Logger) Fatalf(format string, args ...interface{}) { l.Logf(LevelFatal, format, args...) }

func New() *Logger {
	return NewWith(log.Default())
}

func NewWith(delegate Printer) *Logger {
	return &Logger{instance: delegate}
}

var loggerInstance = struct {
	once     sync.Once
	instance *Logger
}{}

func Default() *Logger {
	loggerInstance.once.Do(func() {
		loggerInstance.instance = New()
	})
	return loggerInstance.instance
}
