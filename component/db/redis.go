// redis 包提供redis客户端的连接与初始化
package db

import (
	"github.com/go-redis/redis"
	"github.com/ebar-go/ego/helper"
)


// ConnectRedis 连接服务器
func  ConnectRedis(option *redis.Options) (*redis.Client, error){
	client := redis.NewClient(option)

	pong, err := client.Ping().Result()
	helper.Debug(pong)

	return client, err
}

