package config

import (
	"github.com/ebar-go/ego/component/etcd"
	"github.com/ebar-go/ego/component/mysql"
	"github.com/ebar-go/ego/component/redis"
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
	pprofKey = "server.pprof"
	swaggerKey = "server.swagger"
	taskKey = "server.task"


	mysqlDsnKey 			  = "mysql.dsn"
	mysqlMaxIdleConnectionsKey = "mysql.maxIdleConnections"
	mysqlMaxOpenConnectionsKey = "mysql.maxOpenConnections"
	mysqlMaxLifeTimeKey = "mysql.maxLifeTime"

	redisHostKey        = "redis.host"
	redisPortKey        = "redis.port"
	redisPassKey        = "redis.pass"
	redisPoolSizeKey    = "redis.poolSize"
	redisMaxRetriesKey  = "redis.maxRetries"
	redisIdleTimeoutKey = "redis.idleTimeout"
	redisCluster = "redis.cluster"

	etcdEndpoints = "etcd.endpoints"
	etcdTimeout = "etcd.timeout"
)




// Config 配置
type Config struct {
	*viper.Viper
	server *serverConf
	mysql *mysql.Config
	redis *redis.Config
	etcd *etcd.Config
	mu *sync.Mutex
}


// serverConf  服务配置
type serverConf struct {
	// 服务名称
	Name string

	// 服务端口号
	Port int

	// 响应日志最大长度
	MaxResponseLogSize int

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
func New() *Config  {
	conf := new(Config)
	conf.Viper = viper.New()
	conf.mu = new(sync.Mutex)
	conf.setDefault()
	return conf
}


func (conf *Config) setDefault() {
	conf.AutomaticEnv()
	conf.SetDefault(systemNameKey, "app")
	conf.SetDefault(httpPortKey, 8080)
	conf.SetDefault(maxResponseLogSizeKey, 2000)
	conf.SetDefault(logPathKey, "/tmp/app.log")
	conf.SetDefault(traceHeaderKey, "gateway-trace")
	conf.SetDefault(httpRequestTimeoutKey, 3)

	conf.SetDefault(redisHostKey, "127.0.0.1")
	conf.SetDefault(redisPortKey, 6379)
	conf.SetDefault(redisPoolSizeKey, 100)
	conf.SetDefault(redisMaxRetriesKey, 3)
	conf.SetDefault(redisIdleTimeoutKey, 5)
}

// LoadFile 加载配置文件
func (conf *Config) LoadFile(path string) error {
	conf.SetConfigFile(path)

	return conf.ReadInConfig()
}

// Server
func (conf *Config) Server() (*serverConf) {
	if conf.server == nil {
		// 加锁防止并发
		conf.mu.Lock()
		defer conf.mu.Unlock()
		conf.server = &serverConf{
			Name:               conf.GetString(systemNameKey),
			Port:               conf.GetInt(httpPortKey),
			MaxResponseLogSize: conf.GetInt(maxResponseLogSizeKey),
			LogPath:            conf.GetString(logPathKey),
			JwtSignKey:         []byte(conf.GetString(jwtSignKey)),
			TraceHeader:        conf.GetString(traceHeaderKey),
			HttpRequestTimeOut: conf.GetInt(httpRequestTimeoutKey),
			Debug:              conf.GetBool(debugKey),
			Pprof: conf.GetBool(pprofKey),
			Swagger: conf.GetBool(swaggerKey),
			Task: conf.GetBool(taskKey),
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
func (conf *Config) Redis() (*redis.Config){
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
			Cluster: conf.GetString(redisCluster),
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