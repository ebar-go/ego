package mysql

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"time"
)

type Database interface {
	GetInstance() *gorm.DB
}

type database struct {
	db *gorm.DB
}

func (d *database) GetInstance() *gorm.DB {
	return d.db
}

func Connect(conf *Config) (Database, error)  {
	sqlDB, err := sql.Open("mysql", conf.Dsn)
	if err != nil {
		return nil, err
	}

	// set pool config
	sqlDB.SetMaxIdleConns(conf.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(conf.MaxLifeTime))

	conn, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return nil, err
	}
	log.Printf("connect database success\n")

	db := &database{db: conn}
	return db, nil
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