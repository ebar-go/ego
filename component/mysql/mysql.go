package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"os"
	"fmt"
	"net"
	"strconv"
)

const (
	defaultIdleConnections = 10
	defaultOpenConnections = 40
)

var group *ConnectionGroup

func init() {
	group = &ConnectionGroup{
		lock:        &sync.Mutex{},
		connections: make(map[string]*gorm.DB),
	}
}

// GetConnectionGroup 获取数据库连接池组
func GetConnectionGroup() *ConnectionGroup {
	return group
}

// ConnectionGroup 数据库连接组
type ConnectionGroup struct {
	lock        *sync.Mutex
	defaultKey  string
	connections map[string]*gorm.DB
}

// Conf 数据库配置
type Conf struct {
	// 数据库标识
	Key string

	// 数据库名称
	Name string

	// 地址
	Host string

	// 端口号
	Port int

	// 用户名
	User string

	// 密码
	Password string

	// 是否为默认连接
	Default bool

	// 连接失败的处理方式
	ConnectFailedHandler func(err error)

	// 是否日志模式，默认不开
	LogMode bool

	// 最大操作连接数
	MaxIdleConnections int

	// 最大打开连接数
	MaxOpenConnections int
}

func CloseConnectionGroup() {
	if len(group.connections) == 0 {
		return
	}

	for _, conn := range group.connections {
		conn.Close()
	}

	return
}

// complete 补全
func (conf *Conf) complete() {
	if conf.MaxIdleConnections == 0 {
		conf.MaxIdleConnections = defaultIdleConnections
	}

	if conf.MaxOpenConnections == 0 {
		conf.MaxOpenConnections = defaultOpenConnections
	}

	if conf.ConnectFailedHandler == nil {
		conf.ConnectFailedHandler = func(err error) {
			fmt.Printf("database %s connect failed:%s", conf.Key, err.Error())
			os.Exit(-1)
		}
	}
}

// GetDsn 获取dsn
func (conf Conf) GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User,
		conf.Password,
		net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port)),
		conf.Name)
}

// 初始化数据库连接池
func InitPool(confItems ...Conf) (err error) {
	// 加锁
	group.lock.Lock()

	defaultConnectionKey := ""
	for key, conf := range confItems {
		// 如果没有设置default选项，则默认取第一个
		if key == 0 {
			defaultConnectionKey = conf.Key
		}

		if conf.Default {
			defaultConnectionKey = conf.Key
		}

		conf.complete()

		connection, err := gorm.Open("mysql", conf.GetDsn())
		if err != nil {
			conf.ConnectFailedHandler(err)
		}

		// 设置是否打印日志
		connection.LogMode(conf.LogMode)
		// 设置连接池
		connection.DB().SetMaxIdleConns(conf.MaxIdleConnections)
		connection.DB().SetMaxOpenConns(conf.MaxOpenConnections)

		group.connections[conf.Key] = connection
	}

	group.defaultKey = defaultConnectionKey

	return nil
}

// 获取数据库连接
func GetConnection() *gorm.DB {
	return GetConnectionByName(group.defaultKey)
}

// GetConnectionByName 根据名称获取数据库连接
func GetConnectionByName(name string) *gorm.DB {
	return group.connections[name]
}
