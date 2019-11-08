package db

import (
	"testing"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

// TestConnectRedis 测试连接
func TestConnectRedis(t *testing.T) {

	var err error

	client, err := ConnectRedis(&redis.Options{
		Addr: "192.168.0.222:6379",
		Password:"",
		DB:0,
	})

	assert.Nil(t, err)


	err = client.Set("key", "value", 0).Err()
	assert.Nil(t, err)

	val, err := client.Get("key").Result()
	assert.Nil(t, err)
	fmt.Println("key", val)

}

