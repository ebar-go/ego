package mysql

import (
	"fmt"
	"net"
	"strconv"
)

// Config Mysql的配置项
type Config struct {
	// 表前缀
	TablePrefix string `mapstructure:"tablePrefix"`

	// data sources
	DataSources []DataSource `mapstructure:"dsn"`

	// 最大空闲连接数
	MaxIdleConnections int `mapstructure:"maxIdleConnections"`

	// 最大打开连接数
	MaxOpenConnections int `mapstructure:"maxOpenConnections"`

	// 最长活跃时间
	MaxLifeTime int `mapstructure:"maxLifeTime"`
}

// DataSource 连接配置
type DataSource struct {
	// host
	Host string `mapstruture:"host"`
	// 端口号
	Port int `mapstructure:"port"`
	// 用户名
	User string `mapstructure:"user"`
	// 密码
	Password string `mapstructure:"password"`
	// 数据库名称
	Name string `mapstructure:"name"`
}


// DsnItems 获取全部dsn资源
func (conf Config) DsnItems() []string {
	var items []string
	for _, dsn := range conf.DataSources {
		items = append(items, dsn.dsn())
	}
	return items
}


// Dsn return mysql dsn
func (conf DataSource) dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User,
		conf.Password,
		net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port)),
		conf.Name)
}