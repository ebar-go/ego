package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"net"
	"strconv"
	"github.com/ebar-go/ego/helper"
)

const (
	defaultIdleConnections = 10
	defaultOpenConnections = 40
	defaultPort = 3306
)

// Open 创建连接池
// 仅提供创建连接池的方法，使用者自己维护连接
func Open(conf Conf) (*gorm.DB, error) {
	// 补全配置项
	conf.complete()

	connection, err := gorm.Open("mysql", conf.dsn())
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

// Conf 配置项
type Conf struct {
	// 数据库名称
	Name string

	// 地址
	Host string

	// 端口号,默认是3306
	Port int

	// 用户名
	User string

	// 密码
	Password string

	// 是否日志模式，默认不开
	LogMode bool

	// 最大操作连接数
	MaxIdleConnections int

	// 最大打开连接数
	MaxOpenConnections int
}

// complete 补全配置
func (conf *Conf) complete() {
	conf.MaxIdleConnections = helper.DefaultInt(conf.MaxIdleConnections, defaultIdleConnections)
	conf.MaxOpenConnections = helper.DefaultInt(conf.MaxOpenConnections, defaultOpenConnections)
	conf.Port = helper.DefaultInt(conf.Port, defaultPort)
}

// dsn 获取dsn
func (conf Conf) dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User,
		conf.Password,
		net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port)),
		conf.Name)
}
