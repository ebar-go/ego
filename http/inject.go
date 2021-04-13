/**
 * @Author: Hongker
 * @Description:
 * @File:  inject
 * @Version: 1.0.0
 * @Date: 2021/4/3 19:32
 */

package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func Inject(container *dig.Container) {
	_ = container.Provide(NewServer)
	_ = container.Provide(func(server *Server) *gin.Engine {
		return server.router
	})
}
