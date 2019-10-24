// redis 包提供redis客户端的连接与初始化
package db

import (
	"github.com/go-redis/redis"
	"fmt"
)


// ConnectRedis 连接服务器
func  ConnectRedis(option *redis.Options) (*redis.Client, error){
	client := redis.NewClient(option)

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	return client, err
}

