package mysql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getConf() Conf {
	return Conf{
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
