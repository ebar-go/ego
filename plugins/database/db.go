package database

import (
	"github.com/ebar-go/ego/utils/structure"
	"gorm.io/gorm"
)

var singleton = structure.NewSingleton(dBPluginConstructor)

type DBPlugin struct {
	db *gorm.DB
}

func dBPluginConstructor() *DBPlugin {
	return &DBPlugin{}
}

// Register 注入一个gorm连接
func Register(db *gorm.DB) {
	singleton.Get().db = db
}

// Instance 使用链接
func Instance() *gorm.DB {
	return singleton.Get().db
}
