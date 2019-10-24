package db

import (
	"testing"
	"fmt"
	"github.com/ebar-go/ego/test"
	"github.com/go-redis/redis"
)

// TestConnectRedis 测试连接
func TestConnectRedis(t *testing.T) {

	var err error

	client, err := ConnectRedis(&redis.Options{
		Addr: "192.168.0.222:6379",
		Password:"",
		DB:0,
	})

	test.AssertNil(t, err)


	err = client.Set("key", "value", 0).Err()
	test.AssertNil(t, err)

	val, err := client.Get("key").Result()
	test.AssertNil(t, err)
	fmt.Println("key", val)

}

