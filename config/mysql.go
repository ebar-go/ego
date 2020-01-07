package config

import (
	"fmt"
	"github.com/ebar-go/ego/utils/number"
	"github.com/ebar-go/ego/utils/strings"
	"net"
	"strconv"
)

const (
	mysqlDefaultIdleConnections = 10
	mysqlDefaultOpenConnections = 40
	mysqlDefaultPort            = 3306
	mysqlDefaultUser = "root"
)

// MysqlOptions
type MysqlOptions struct {
	// auto connect
	AutoConnect bool

	// database name
	Name string

	// host
	Host string

	// port, default 3306
	Port int

	// user, default root
	User string

	// password
	Password string

	// log mode
	LogMode bool

	// MaxIdleConnections, default 10
	MaxIdleConnections int

	// MaxOpenConnections, default 40
	MaxOpenConnections int
}

// complete set default options
func (options *MysqlOptions) complete() {
	options.MaxIdleConnections = number.DefaultInt(options.MaxIdleConnections, mysqlDefaultIdleConnections)
	options.MaxOpenConnections = number.DefaultInt(options.MaxOpenConnections, mysqlDefaultOpenConnections)
	options.Port = number.DefaultInt(options.Port, mysqlDefaultPort)
	options.User = strings.Default(options.User, mysqlDefaultUser)
}

// Dsn return mysql dsn
func (options MysqlOptions) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		options.User,
		options.Password,
		net.JoinHostPort(options.Host, strconv.Itoa(options.Port)),
		options.Name)
}
