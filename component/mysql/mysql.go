package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"os"
	"fmt"
)

const (
	defaultIdleConnections = 10
	defaultOpenConnections = 40
)

var group *ConnectionGroup

func init()  {
	group = &ConnectionGroup{
		lock: &sync.Mutex{},
		connections: make(map[string]*gorm.DB),
	}
}

// ConnectionGroup 数据库连接组
type ConnectionGroup struct {
	lock *sync.Mutex
	defaultName string
	connections map[string]*gorm.DB
}

// Conf 数据库配置
type Conf struct {
	// 数据库名称
	Name string

	// 是否为默认连接
	Default bool

	// 连接失败的处理方式
	ConnectFailedHandler func(err error)

	// dsn
	Dsn string

	// 是否日志模式，默认不开
	LogMode bool

	// 最大操作连接数
	MaxIdleConnections int

	// 最大打开连接数
	MaxOpenConnections int
}

func CloseConnectionGroup()  {
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
			fmt.Printf("database %s connect failed:%s", conf.Name, err.Error())
			os.Exit(-1)
		}
	}
}

// 初始化数据库连接池
func InitPool(confItems ...Conf) (err error) {
	group.lock.Lock()

	defaultConnectionName := ""
	for key, conf := range confItems {
		// 如果没有设置default选项，则默认取第一个
		if key == 0 {
			defaultConnectionName = conf.Name
		}else if conf.Default {
			defaultConnectionName = conf.Name
		}

		conf.complete()

		connection, err := gorm.Open("mysql", conf.Dsn)
		if err != nil {
			conf.ConnectFailedHandler(err)
		}

		// 设置是否打印日志
		connection.LogMode(conf.LogMode)
		// 设置连接池
		connection.DB().SetMaxIdleConns(conf.MaxIdleConnections)
		connection.DB().SetMaxOpenConns(conf.MaxOpenConnections)

		group.connections[conf.Name] = connection
	}
	group.defaultName = defaultConnectionName

	return nil
}

// 获取数据库连接
func GetConnection() *gorm.DB {
	return GetConnectionByName(group.defaultName)
}

// GetConnectionByName 根据名称获取数据库连接
func GetConnectionByName(name string) *gorm.DB {
	return group.connections[name]
}