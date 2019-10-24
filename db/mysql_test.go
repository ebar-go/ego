package db

import (
	"testing"
	"git.epetbar.com/go-package/ego/test"
	"fmt"
)


func TestConnectMysql(t *testing.T) {
	client, err := ConnectMysql("root:123456@tcp(192.168.0.212:3306)/epet_stock?charset=utf8&parseTime=True&loc=Local")
	fmt.Println(err)
	test.AssertNil(t, err)

	// 启用Logger，显示详细日志
	client.LogMode(true)

	// 禁用日志记录器，不显示任何日志
	client.LogMode(false)

	client.DB().Ping()

	defer client.Close()
}
