package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/ego/http/constant"
	"github.com/ebar-go/ego/log"
	"github.com/ebar-go/ego/http/helper"
)

// Recover
func Recover(ctx *gin.Context)  {
	defer func() {
		if r := recover(); r != nil {
			response.Error(ctx, constant.StatusError, "系统错误")

			log.System().Error("system_error", log.Context{
				"trace_id" : helper.GetTraceId(ctx),
				"error" : r,
			})
		}
	}()
	ctx.Next()
}
