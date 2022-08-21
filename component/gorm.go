package component

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Gorm struct {
	Named
	*gorm.DB
}

// Open connect to database
func (g *Gorm) Open(dialector gorm.Dialector, config *gorm.Config) error {
	conn, err := gorm.Open(dialector, config)
	if err != nil {
		g.DB = conn
	}

	return err
}

// OpenMySQL connect to mysql data
func (g *Gorm) OpenMySQL(dsn string) error {
	return g.Open(mysql.Open(dsn), &gorm.Config{})
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
