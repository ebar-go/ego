package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"fmt"
	"os"
)

const (
	defaultIdleConnections = 10
	defaultOpenConnections = 40
)

// Conf 数据库配置
type Conf struct {
	// dsn
	Dsn string

	// 是否日志模式，默认不开
	LogMode bool

	// 最大操作连接数
	MaxIdleConnections int

	// 最大打开连接数
	MaxOpenConnections int
}

var dbInitOnce sync.Once
var connection *gorm.DB

// complete 补全
func (conf *Conf) complete() {
	if conf.MaxIdleConnections == 0 {
		conf.MaxIdleConnections = defaultIdleConnections
	}

	if conf.MaxOpenConnections == 0 {
		conf.MaxOpenConnections = defaultOpenConnections
	}
}

// 初始化数据库连接池
func InitPool(conf Conf) (err error) {

	dbInitOnce.Do(func() {
		conf.complete()

		client, err := gorm.Open("mysql", conf.Dsn)
		if err != nil {
			fmt.Println("数据库连接失败:", err)
			os.Exit(-1)
		}

		connection = client

		// 设置是否打印日志
		connection.LogMode(conf.LogMode)
		// 设置连接池
		connection.DB().SetMaxIdleConns(conf.MaxIdleConnections)
		connection.DB().SetMaxOpenConns(conf.MaxOpenConnections)
	})

	return nil
}

// 获取数据库连接
func GetConnection() *gorm.DB {
	return connection
}