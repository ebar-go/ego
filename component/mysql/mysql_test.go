package mysql

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func getConf() Conf {
	return Conf{
		Key: "db0",
		Host: "127.0.0.1",
		Port: 3306,
		User: "root",
		Password: "123456",
		Default: true,
		LogMode:true,
	}
}

func TestInitPool(t *testing.T) {
	conf := getConf()
	err := InitPool(conf)
	assert.Nil(t, err)

	defer GetConnection().Close()
}

func TestGetConnection(t *testing.T) {
	conf := getConf()
	err := InitPool(conf)
	assert.Nil(t, err)

	conn := GetConnection()
	pingErr := conn.DB().Ping()
	assert.Nil(t, pingErr)

	defer conn.Close()
}
