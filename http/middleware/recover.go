package middleware

import (
	"fmt"
	"github.com/ebar-go/ego/helper"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/ego/log"
	"github.com/gin-gonic/gin"
)

// Recover
func Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			response.Error(ctx, 500, "系统错误")
			fmt.Println(helper.Trace())

			log.System().Error("system_error", log.Context{
				"error": r,
			})
		}
	}()
	ctx.Next()
}
