package handler

import (
	"github.com/ebar-go/ego/app"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	// 禁用swagger的环境变量标识
	disableEnv = "DisableSwagger"
)

// Swagger
func SwaggerHandler() gin.HandlerFunc {
	if app.Config().Server().Swagger {
		return ginSwagger.WrapHandler(swaggerFiles.Handler)
	}
	return func(c *gin.Context) {
		// Simulate behavior when route unspecified and
		// return 404 HTTP code
		c.String(404, "")
	}
}
