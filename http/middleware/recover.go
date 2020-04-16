package middleware

import (
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/ego/utils"
	"github.com/gin-gonic/gin"
)

// Recover
func Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			response.WrapContext(ctx).Error(500, "System Error")

			log.System().Error("system_error", log.Context(map[string]interface{}{
				"error": r,
				"trace": utils.Trace(),
			}))
		}
	}()
	ctx.Next()
}
