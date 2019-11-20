package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CORS 跨域中间件
func CORS(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin")

	// 核心处理方式
	c.Header("Access-Control-Allow-Origin", origin)
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")

	if method == "OPTIONS" || method == "HEAD" {
		c.JSON(http.StatusNoContent, "")
		c.Abort()
		return
	}

	c.Next()

}
