package app

import (
	"github.com/ebar-go/ego/component/config"
	"github.com/ebar-go/ego/component/log"
)

func newLogger(conf *config.Config) *log.Logger {
	return log.New(conf.Debug, map[string]interface{}{"app_name":conf.Name})
}
