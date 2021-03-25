package server

import (
	"github.com/ebar-go/ego/component/config"
	"go.uber.org/dig"
)
const (
	envKey     = "server.environment"
	systemNameKey         = "server.systemName"
	httpPortKey           = "server.httpPort"
	maxResponseLogSizeKey = "server.maxResponseLogSize"
	logPathKey            = "server.logPath"
	traceHeaderKey        = "server.traceHeader"
	httpRequestTimeoutKey = "server.httpRequestTimeout"
	jwtSignKey            = "server.jwtSign"
	debugKey              = "server.debug"
	pprofKey              = "server.pprof"
	swaggerKey            = "server.swagger"
	taskKey               = "server.task"
)
func NewConfig(conf *config.Config) *Config {
	return &Config{
		Environment:        conf.GetString(envKey),
		Name:               conf.GetString(systemNameKey),
		Port:               conf.GetInt(httpPortKey),
		MaxResponseLogSize: conf.GetInt(maxResponseLogSizeKey),
		LogPath:            conf.GetString(logPathKey),
		JwtSignKey:         []byte(conf.GetString(jwtSignKey)),
		TraceHeader:        conf.GetString(traceHeaderKey),
		HttpRequestTimeOut: conf.GetInt(httpRequestTimeoutKey),
		Debug:              conf.GetBool(debugKey),
		Pprof:              conf.GetBool(pprofKey),
		Swagger:            conf.GetBool(swaggerKey),
		Task:               conf.GetBool(taskKey),
	}
}

func Inject(container *dig.Container) error {
	return container.Provide(NewConfig)
}
