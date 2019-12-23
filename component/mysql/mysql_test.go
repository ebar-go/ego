package mysql

import (
	"github.com/ebar-go/ego/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getConf() config.MysqlConfig {
	return config.MysqlConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "123456",
		LogMode:  true,
	}
}

func TestOpen(t *testing.T) {
	conf := getConf()
	conn, err := Open(conf)
	assert.Nil(t, err)

	defer conn.Close()
}
