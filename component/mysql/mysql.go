package mysql

import (
	"github.com/ebar-go/ego/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Open 创建连接池
func Open(dsn string) (*gorm.DB, error) {

	return gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 设置是否打印日志
	connection.LogMode(conf.LogMode)
	// 设置连接池
	connection.DB().SetMaxIdleConns(conf.MaxIdleConnections)
	connection.DB().SetMaxOpenConns(conf.MaxOpenConnections)

	return connection, nil
}
