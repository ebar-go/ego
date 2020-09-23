package mysql

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

type Manager struct {
	*gorm.DB

	conf *Config
}

func NewManager(conf *Config) *Manager {
	return &Manager{conf: conf}
}

func (db *Manager) Connect() error {
	sqlDB, err := sql.Open("mysql", db.conf.Dsn)
	if err != nil {
		return err
	}

	// set pool config
	sqlDB.SetMaxIdleConns(db.conf.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(db.conf.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(db.conf.MaxLifeTime))

	conn, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return err
	}
	fmt.Printf("connect database success\n")

	db.DB = conn

	return nil
}

func Resolver() *dbresolver.DBResolver {
	return new(dbresolver.DBResolver)
}

func ResolverConfig(item ResolverItem) dbresolver.Config {
	var sources, replicas []gorm.Dialector
	for _, source := range item.Sources {
		sources = append(sources, mysql.Open(source))
	}

	for _, replica := range item.Replicas {
		replicas = append(replicas, mysql.Open(replica))
	}

	return dbresolver.Config{
		Sources:  sources,
		Replicas: replicas,
	}
}
