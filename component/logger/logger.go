package logger

import (
	"fmt"
	"log"
)

type Printer interface {
	Println(v ...any)
}

type Instance struct {
	printer Printer
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

func (l *Instance) Log(level LogLevel, message string) { l.printer.Println(level, message) }
func (l *Instance) Logf(level LogLevel, format string, a ...interface{}) {
	l.printer.Println(level, fmt.Sprintf(format, a...))
}

func (l *Instance) Info(message string)  { l.Log(LevelInfo, message) }
func (l *Instance) Debug(message string) { l.Log(LevelDebug, message) }
func (l *Instance) Warn(message string)  { l.Log(LevelWarn, message) }
func (l *Instance) Error(message string) { l.Log(LevelError, message) }
func (l *Instance) Fatal(message string) { l.Log(LevelFatal, message) }
func (l *Instance) Panic(message string) { l.Log(LevelPanic, message) }

func (l *Instance) Infof(format string, args ...interface{})  { l.Logf(LevelInfo, format, args...) }
func (l *Instance) Debugf(format string, args ...interface{}) { l.Logf(LevelDebug, format, args...) }
func (l *Instance) Warnf(format string, args ...interface{})  { l.Logf(LevelWarn, format, args...) }
func (l *Instance) Errorf(format string, args ...interface{}) { l.Logf(LevelError, format, args...) }
func (l *Instance) Fatalf(format string, args ...interface{}) { l.Logf(LevelFatal, format, args...) }

func (l *Instance) SetPrinter(printer Printer) {
	l.printer = printer
}

func New() *Instance {
	return &Instance{printer: log.Default()}
}
