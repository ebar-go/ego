package log

// init system log
// App : use for application logic
// Request : use for http request
// System: use for system debug
// Mq: use for mq

import (
	"github.com/ebar-go/ego/config"
	"go.uber.org/zap"
	"sync"
)

var appLogger *Logger
var mqLogger *Logger
var systemLogger *Logger
var requestLogger *Logger

// Logger
type Logger struct {
	once sync.Once
	instance *zap.Logger
}

func init()  {
	appLogger = new(Logger)
	mqLogger = new(Logger)
	systemLogger = new(Logger)
	requestLogger = new(Logger)
}

// getInstance init logger instance
func (l *Logger) getInstance(filename string) *zap.Logger {
	serverConfig := config.Server()

	// init once
	l.once.Do(func() {
		l.instance = NewZap(serverConfig.LogPath + "/" + filename, zap.InfoLevel,
			zap.String("system_name", serverConfig.Name), zap.Int("system_port", serverConfig.Port))
	})

	return l.instance
}

// MQ return mns log
func MQ() *zap.Logger {
	return mqLogger.getInstance("mq.log")
}

// System return system log
func System() *zap.Logger {
	return systemLogger.getInstance("system.log")
}

// Request return request log
func Request() *zap.Logger {
	return requestLogger.getInstance("request.log")
}

func App() *zap.Logger {
	return appLogger.getInstance("app.log")
}

