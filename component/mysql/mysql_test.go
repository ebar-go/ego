package mysql

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func getConf() Conf {
	return Conf{
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

func TestManyArg(t *testing.T)  {
	many(4,5,6)
}

func many(ids ...int)  {
	for key, value := range ids {
		fmt.Println(key, value)
	}
}
