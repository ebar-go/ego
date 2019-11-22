package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS(c *gin.Context) {
	method := c.Request.Method

	// 核心处理方式
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")

	if method == "OPTIONS" || method == "HEAD" {
		c.AbortWithStatus(204)

		return
	}

	c.Next()

}
