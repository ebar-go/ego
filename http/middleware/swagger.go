package middleware

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	// 禁用swagger的环境变量标识
	disableEnv = "DisableSwagger"
)

// Swagger
func Swagger() gin.HandlerFunc {
	return ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, disableEnv)
}
