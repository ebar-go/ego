package gorm

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

type Instance struct {
	*gorm.DB
}

func New() *Instance {
	return &Instance{}
}

// Connect open mysql connection
func (instance *Instance) Connect(dsn string, config *gorm.Config) (err error) {
	instance.DB, err = gorm.Open(mysql.Open(dsn), config)
	return
}

// RegisterResolverConfig registers resolver configuration for current connection
func RegisterResolverConfig(db *gorm.DB, config dbresolver.Config, tables ...interface{}) error {
	resolver := &dbresolver.DBResolver{}
	// get resolver plugin from config
	plugin, ok := db.Config.Plugins[resolver.Name()]
	if !ok {
		// if resolver not exist, create and initialize it.
		resolver = dbresolver.Register(config, tables...)
		return db.Use(resolver)
	}
	// if resolver is already exist, register configuration directly.
	// Because of this plugin is a pointer, it will affect when use register.
	plugin.(*dbresolver.DBResolver).Register(config, tables...)
	return nil
}

// EnableConnectionPool enables connection pool
func EnableConnectionPool(db *sql.DB, maxIdleConns int, maxOpenConns int, connMaxLifetime time.Duration) {
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	db.SetMaxIdleConns(maxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	db.SetMaxOpenConns(maxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	db.SetConnMaxLifetime(connMaxLifetime)
}
