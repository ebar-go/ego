package config

import (
	"github.com/spf13/viper"
)

const (
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
	envKey                = "server.environment"
)

// Config 配置
type Config struct {
	*viper.Viper
	// 运行环境
	Environment string
	// 服务名称
	Name string
	// 服务端口号,
	Port int
	// 日志路径
	LogPath string
	// jwt的key
	JwtSignKey []byte
	// trace header key
	TraceHeader string
	// http request timeout
	HttpRequestTimeOut int
	// 是否开启debug,开启后会显示debug信息
	Debug bool
	// 是否开启pprof
	Pprof bool
	// 是否开启swagger文档
	Swagger bool
	// 是否开启定时任务
	Task bool
}

// New 实例
func New() *Config {
	conf := new(Config)
	conf.Viper = viper.New()
	conf.setDefaults()
	return conf
}

func (conf *Config) setDefaults() {
	conf.SetDefault(systemNameKey, "app")
	conf.SetDefault(httpPortKey, 8080)
	conf.SetDefault(maxResponseLogSizeKey, 2000)
	conf.SetDefault(logPathKey, "/tmp/app.log")
	conf.SetDefault(traceHeaderKey, "gateway-trace")
	conf.SetDefault(httpRequestTimeoutKey, 3)
}

// LoadFile 加载配置文件
func (conf *Config) LoadFile(path ...string) error {
	for _, p := range path {
		conf.SetConfigFile(p)
		if err := conf.MergeInConfig(); err != nil {
			return err
		}
	}

	conf.Environment = conf.GetString(envKey)
	conf.Name = conf.GetString(systemNameKey)
	conf.Port = conf.GetInt(httpPortKey)
	conf.LogPath = conf.GetString(logPathKey)
	conf.JwtSignKey = []byte(conf.GetString(jwtSignKey))
	conf.TraceHeader = conf.GetString(traceHeaderKey)
	conf.HttpRequestTimeOut = conf.GetInt(httpRequestTimeoutKey)
	conf.Debug = conf.GetBool(debugKey)
	conf.Pprof = conf.GetBool(pprofKey)
	conf.Swagger = conf.GetBool(swaggerKey)
	conf.Task = conf.GetBool(taskKey)

	return nil
}

const (
	envProduct = "product"
	envDevelop = "develop"
)

// IsProduct 是否为生产环境
func (conf *Config) IsProduct() bool {
	return envProduct == conf.Environment
}

// IsDevelop 是否为测试环境
func (conf *Config) IsDevelop() bool {
	return envDevelop == conf.Environment
}
