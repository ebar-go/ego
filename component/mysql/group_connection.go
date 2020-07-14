package mysql

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

// GroupManager manage mysql connections
type GroupManager struct {
	connections map[string]*gorm.DB
	options map[string]Config
	dialect string
}

// New
func New(options map[string]Config) *GroupManager  {
	return &GroupManager{connections: make(map[string]*gorm.DB), dialect:"mysql", options:options}
}

// GetConnection 获取连接
func (gm *GroupManager) GetConnection(name string) *gorm.DB  {
	return gm.connections[name]
}

func (gm *GroupManager) Connect() error {
	for name, conf := range gm.options {
		adapter, err := NewReadWriteAdapter(gm.dialect, conf.DsnItems())
		if err != nil {
			return err
		}

		adapter.SetMaxIdleConns(conf.MaxIdleConnections)
		adapter.SetMaxOpenConns(conf.MaxOpenConnections)
		adapter.SetConnMaxLifetime(time.Duration(conf.MaxLifeTime) * time.Second)

		conn, err := gorm.Open(gm.dialect, adapter)
		if err != nil {
			return err
		}

		if err := adapter.Ping(); err != nil {
			return err
		}
		log.Println("Connect mysql success:", name)
		gm.connections[name] = conn
	}


	return nil
}
