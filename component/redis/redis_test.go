package redis

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func getConf() Conf {
	return Conf{
		Host: "192.168.0.222",
		Port: 6379,
	}
}

// TestInitPool 测试初始化连接池
func TestInitPool(t *testing.T) {

	var err error
	err = InitPool(getConf())

	assert.Nil(t, err)

	defer GetConnection().Close()

}

func TestGetConnection(t *testing.T) {
	var err error
	err = InitPool(getConf())

	assert.Nil(t, err)

	client := GetConnection()
	defer client.Close()
	err = client.Set("key", "value", 0).Err()
	assert.Nil(t, err)

	val, err := client.Get("key").Result()
	assert.Nil(t, err)
	fmt.Println("key", val)
}

