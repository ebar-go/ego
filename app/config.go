package app

import (
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/etcd"
	"github.com/ebar-go/ego/component/mysql"
	"github.com/ebar-go/ego/component/redis"
	"github.com/ebar-go/ego/http"
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
	envKey     = "server.environment"
)
const (
	mysqlDsnKey                = "mysql.dsn"
	mysqlMaxIdleConnectionsKey = "mysql.maxIdleConnections"
	mysqlMaxOpenConnectionsKey = "mysql.maxOpenConnections"
	mysqlMaxLifeTimeKey        = "mysql.maxLifeTime"
)
const (
	redisHostKey        = "redis.host"
	redisPortKey        = "redis.port"
	redisPassKey        = "redis.pass"
	redisPoolSizeKey    = "redis.poolSize"
	redisMaxRetriesKey  = "redis.maxRetries"
	redisIdleTimeoutKey = "redis.idleTimeout"
	redisCluster        = "redis.cluster"
)
const (
	etcdEndpoints = "etcd.endpoints"
	etcdTimeout   = "etcd.timeout"
)

func newHttpConfig(conf *config.Config) *http.Config {
	conf.SetDefault(systemNameKey, "app")
	conf.SetDefault(httpPortKey, 8080)
	conf.SetDefault(maxResponseLogSizeKey, 2000)
	conf.SetDefault(logPathKey, "/tmp/app.log")
	conf.SetDefault(traceHeaderKey, "gateway-trace")
	conf.SetDefault(httpRequestTimeoutKey, 3)

	return &http.Config{
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

func newDatabaseConfig(conf *config.Config) *mysql.Config {
	return &mysql.Config{
		MaxIdleConnections: conf.GetInt(mysqlMaxIdleConnectionsKey),
		MaxOpenConnections: conf.GetInt(mysqlMaxOpenConnectionsKey),
		MaxLifeTime:        conf.GetInt(mysqlMaxLifeTimeKey),
		Dsn:                conf.GetString(mysqlDsnKey),
	}
}

func newRedisConfig(conf *config.Config) *redis.Config {
	return &redis.Config{
		Host:        conf.GetString(redisHostKey),
		Port:        conf.GetInt(redisPortKey),
		Auth:        conf.GetString(redisPassKey),
		PoolSize:    conf.GetInt(redisPoolSizeKey),
		MaxRetries:  conf.GetInt(redisMaxRetriesKey),
		IdleTimeout: time.Duration(conf.GetInt(redisIdleTimeoutKey)) * time.Second,
		Cluster:     conf.GetStringSlice(redisCluster),
	}
}

func newEtcdConfig(conf *config.Config) *etcd.Config {
	return &etcd.Config{
		Endpoints: conf.GetStringSlice(etcdEndpoints),
		Timeout:   conf.GetInt(etcdTimeout),
	}
}