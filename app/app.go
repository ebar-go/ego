package app

import (
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ego/container"
	"github.com/ebar-go/ego/helper"
	"github.com/ebar-go/ego/log"
)

var (
	app = container.New()
)

// Config 配置项
func Config() *config.Config {
	var conf *config.Config

	if err := app.Invoke(func(c *config.Config) {
		conf = c
	}); err != nil {
		helper.CheckError("InitConfig", app.Provide(config.NewInstance))
		conf = config.NewInstance()
	}

	return conf
}

func LogManager() log.Manager {
	var manager log.Manager
	if err := app.Invoke(func(m log.Manager) {
		manager = m
	}); err != nil {
		conf := Config()
		helper.CheckError("InitConfig", app.Provide(func() log.Manager{
			m := log.NewManager(log.ManagerConf{
				SystemName: conf.ServiceName,
				SystemPort: conf.ServicePort,
				LogPath:    conf.LogPath,
			})

			manager = m
			return m
		}))
	}

	return manager
}

