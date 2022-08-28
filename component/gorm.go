package component

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

type Gorm struct {
	Named
	*gorm.DB
}

// Open connect to database
func (g *Gorm) Open(dialector gorm.Dialector, config *gorm.Config) (err error) {
	g.DB, err = gorm.Open(dialector, config)
	return
}

// OpenMySQL connect to mysql data
func (g *Gorm) OpenMySQL(dsn string) error {
	return g.Open(mysql.Open(dsn), &gorm.Config{})
}

// UseResolver use resolver for gorm component.include set many db connection configuration with different model or table name
func (g *Gorm) UseResolver(resolver *dbresolver.DBResolver) error {
	return g.Use(resolver)
}

// RegisterResolverConfig registers resolver configuration for current connection
func (g *Gorm) RegisterResolverConfig(config dbresolver.Config, tables ...interface{}) {
	resolver := &dbresolver.DBResolver{}
	// get resolver plugin from config
	plugin, ok := g.Config.Plugins[resolver.Name()]
	if !ok {
		// if resolver not exist, create and initialize it.
		resolver = dbresolver.Register(config, tables)
		g.Use(resolver)
		return
	}
	// if resolver is already exist, register configuration directly.Because of this plugin is a pointer, it will effect when use register.
	plugin.(*dbresolver.DBResolver).Register(config, tables)
}

// EnableConnectionPool enables connection pool
func (g *Gorm) EnableConnectionPool(maxIdleConns int, maxOpenConns int, connMaxLifetime time.Duration) error {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(maxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(maxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	return nil
}

func NewGorm() *Gorm {
	return &Gorm{Named: componentGorm}
}
