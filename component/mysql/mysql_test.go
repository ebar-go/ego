package mysql

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func getConf() Conf {
	return Conf{
		Dsn: "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local",
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
