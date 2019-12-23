package redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getConf() Conf {
	return Conf{
		Host: "192.168.0.212",
		Port: 6379,
	}
}

// TestInitPool 测试初始化连接池
func TestInitPool(t *testing.T) {

	cli, err := Open(getConf())

	assert.Nil(t, err)

	defer cli.Close()

}
