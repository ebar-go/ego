package log

import (
	"github.com/ebar-go/ego/component/trace"
	"go.uber.org/zap"
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
func (logger *Logger) Info(message string, ctx Context) {
	logger.instance.Info(message, format(ctx))
}

// Debug
func (logger *Logger) Debug(message string, ctx Context) {
	logger.instance.Debug(message, format(ctx))
}

// Error
func (logger *Logger) Error(message string, ctx Context) {
	logger.instance.Error(message, format(ctx))
}

// Logger
type Logger struct {
	path     string
	debug    bool
	fields   map[string]interface{}
	instance *zap.Logger
}

// getInstance init logger instance
func New(logPath string, debug bool, fields map[string]interface{}) *Logger {
	logger := new(Logger)
	logger.path = logPath
	logger.debug = debug
	logger.fields = fields
	logger.init()

	return logger

}

func (logger *Logger) init() {
	level := zap.InfoLevel
	if logger.debug {
		level = zap.DebugLevel
	}

	var fields []zap.Field
	for idx, val := range logger.fields {
		fields = append(fields, zap.Any(idx, val))
	}
	logger.instance = newZap(logger.path, level, fields...)
}
