package middleware

import "github.com/gin-gonic/gin"

// Favicon filter /favicon.ico
func Favicon(ctx *gin.Context)  {
	if ctx.Request.RequestURI == "/favicon.ico" {
		ctx.AbortWithStatus(204)
		return
	}

	ctx.Next()
}
