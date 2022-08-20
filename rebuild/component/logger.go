package component

import (
	"fmt"
	"log"
)

type Logger struct {
}

const (
	LevelInfo = iota
	LevelDebug
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

func (l *Logger) Log(level int, message string) { log.Println(level, message) }
func (l *Logger) Logf(level int, format string, a ...interface{}) {
	log.Println(level, fmt.Sprintf(format, a...))
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

func NewLogger() *Logger {
	return &Logger{}
}
