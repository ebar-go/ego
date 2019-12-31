package middleware

import (
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/http/response"
	"github.com/gin-gonic/gin"
)

// Recover
func Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			response.WrapperContext(ctx).Error(500, "System Error")

			app.LogManager().System().Error("system_error", log.Context{
				"error": r,
			})
		}
	}()
	ctx.Next()
}
