package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// ConnectMysql 连接Mysql
func ConnectMysql(dsn string) (*gorm.DB, error){
	client, err := gorm.Open("mysql", dsn)
	return client, err
}

