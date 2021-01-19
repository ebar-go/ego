package config

import (
	"github.com/ebar-go/ego/component/etcd"
	"github.com/ebar-go/ego/component/mysql"
	"github.com/ebar-go/ego/component/redis"
	"github.com/ebar-go/ego/http/server"
	"github.com/spf13/viper"
	"sync"
	"time"
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

	mysqlDsnKey                = "mysql.dsn"
	mysqlMaxIdleConnectionsKey = "mysql.maxIdleConnections"
	mysqlMaxOpenConnectionsKey = "mysql.maxOpenConnections"
	mysqlMaxLifeTimeKey        = "mysql.maxLifeTime"

	redisHostKey        = "redis.host"
	redisPortKey        = "redis.port"
	redisPassKey        = "redis.pass"
	redisPoolSizeKey    = "redis.poolSize"
	redisMaxRetriesKey  = "redis.maxRetries"
	redisIdleTimeoutKey = "redis.idleTimeout"
	redisCluster        = "redis.cluster"

	etcdEndpoints = "etcd.endpoints"
	etcdTimeout   = "etcd.timeout"

	envKey     = "server.environment"

)

const(
	envProduct = "product"
	envDevelop = "develop"
)

// Config 配置
type Config struct {
	*viper.Viper
	server *server.Config
	mysql  *mysql.Config
	redis  *redis.Config
	etcd   *etcd.Config
	mu     *sync.Mutex
}


// New 实例
func New() *Config {
	conf := new(Config)
	conf.Viper = viper.New()
	conf.mu = new(sync.Mutex)
	conf.setDefault()
	return conf
}


// LoadFile 加载配置文件
func (conf *Config) LoadFile(path ...string) error {
	for _, p := range path {
		conf.SetConfigFile(p)
		if err := conf.MergeInConfig(); err != nil {
			return err
		}
	}

	return nil
}


func (conf *Config) setDefault() {
	conf.AutomaticEnv()
	conf.SetDefault(systemNameKey, "app")
	conf.SetDefault(httpPortKey, 8080)
	conf.SetDefault(maxResponseLogSizeKey, 2000)
	conf.SetDefault(logPathKey, "/tmp/app.log")
	conf.SetDefault(traceHeaderKey, "gateway-trace")
	conf.SetDefault(httpRequestTimeoutKey, 3)
}


// Server
func (conf *Config) Server() *server.Config {
	if conf.server == nil {
		// 加锁防止并发
		conf.mu.Lock()
		defer conf.mu.Unlock()
		conf.server = &server.Config{
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

	return conf.server
}

// mysql
func (conf *Config) Mysql() *mysql.Config {
	if conf.mysql == nil {
		conf.mu.Lock()
		defer conf.mu.Unlock()
		conf.mysql = &mysql.Config{
			MaxIdleConnections: conf.GetInt(mysqlMaxIdleConnectionsKey),
			MaxOpenConnections: conf.GetInt(mysqlMaxOpenConnectionsKey),
			MaxLifeTime:        conf.GetInt(mysqlMaxLifeTimeKey),
			Dsn:                conf.GetString(mysqlDsnKey),
		}
	}

	return conf.mysql
}

// Redis
func (conf *Config) Redis() *redis.Config {
	if conf.redis == nil {
		conf.mu.Lock()
		defer conf.mu.Unlock()
		conf.redis = &redis.Config{
			Host:        conf.GetString(redisHostKey),
			Port:        conf.GetInt(redisPortKey),
			Auth:        conf.GetString(redisPassKey),
			PoolSize:    conf.GetInt(redisPoolSizeKey),
			MaxRetries:  conf.GetInt(redisMaxRetriesKey),
			IdleTimeout: time.Duration(conf.GetInt(redisIdleTimeoutKey)) * time.Second,
			Cluster:     conf.GetStringSlice(redisCluster),
		}
	}
	return conf.redis

}

// etcd
func (conf *Config) Etcd() *etcd.Config {
	if conf.etcd == nil {
		conf.mu.Lock()
		defer conf.mu.Unlock()
		conf.etcd = &etcd.Config{
			Endpoints: conf.GetStringSlice(etcdEndpoints),
			Timeout:   conf.GetInt(etcdTimeout),
		}
	}
	return conf.etcd
}


// IsProduct 是否为生产环境
func (conf *Config) IsProduct() bool {
	return envProduct == conf.Server().Environment
}

// IsDevelop 是否为测试环境
func (conf *Config) IsDevelop() bool {
	return envDevelop == conf.Server().Environment
}
