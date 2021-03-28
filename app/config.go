package app

import (
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/etcd"
	"github.com/ebar-go/ego/component/mysql"
	"github.com/ebar-go/ego/component/redis"
	"time"
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