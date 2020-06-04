package log

import (
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/ego/config"
	"go.uber.org/zap"
	"sync"
)

type Context map[string]interface{}

// format
func format(ctx Context) zap.Field {
	if _, ok := ctx["trace_id"]; !ok {
		ctx["trace_id"] = trace.Get()
	}

	return zap.Any("context", ctx)
}

// Info
func Info(message string, ctx Context) {
	d.getInstance().Info(message, format(ctx))
}

// Debug
func Debug(message string, ctx Context) {
	d.getInstance().Debug(message, format(ctx))
}

// Error
func Error(message string, ctx Context) {
	d.getInstance().Error(message, format(ctx))
}

var d = new(Logger)

// Logger
type Logger struct {
	once     sync.Once
	instance *zap.Logger
}

// getInstance init logger instance
func (l *Logger) getInstance() *zap.Logger {
	serverConfig := config.Server()

	level := zap.InfoLevel
	if serverConfig.Debug {
		level = zap.DebugLevel
	}
	// init once
	l.once.Do(func() {
		l.instance = newZap(serverConfig.LogPath, level,
			zap.String("system_name", serverConfig.Name),
			zap.Int("system_port", serverConfig.Port))
	})

	return l.instance
}
