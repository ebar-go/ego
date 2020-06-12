package redis

import (
	goredis "github.com/go-redis/redis"
)

type Redis struct {
	Options *goredis.Options
	*goredis.Client
}


func (redis *Redis) Connect() error {
	connection := goredis.NewClient(redis.Options)
	_, err := connection.Ping().Result()
	if err != nil {
		return err
	}
	redis.Client = connection
	return nil
}

type RedisCluster struct {
	Options *goredis.ClusterOptions
	*goredis.ClusterClient
}

func (redis *RedisCluster) Connect() error {
	connection := goredis.NewClusterClient(redis.Options)
	_, err := connection.Ping().Result()
	if err != nil {
		return err
	}
	redis.ClusterClient = connection
	return nil
}
