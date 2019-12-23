package event

import (
	"github.com/ebar-go/ego/config"
	"github.com/ebar-go/ego/container"
	"github.com/ebar-go/ego/helper"
	"github.com/ebar-go/ego/http/middleware"
	"github.com/ebar-go/ego/log"
)

func PrepareContainer() {
	c := container.New()

	// 初始化JWT
	helper.FatalError("InitJwt", c.Provide(func() middleware.Jwt {
		return middleware.NewJwt(config.Instance.JwtSignKey)
	}))

	container.App = c
}

func PrepareLogManager() {
	log.InitManager(log.ManagerConf{
		SystemName: config.Instance.ServiceName,
		SystemPort: config.Instance.ServicePort,
		LogPath:    config.Instance.LogPath,
	})
}
