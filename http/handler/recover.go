package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ebar-go/ego/http/response"
	"github.com/ebar-go/ego/http/constant"
	"github.com/ebar-go/ego/log"
	"github.com/ebar-go/ego/http/helper"
	"github.com/ebar-go/ego/library"
)

func Recover(ctx *gin.Context)  {
	defer func() {
		if r := recover(); r != nil {
			response.Error(ctx, constant.StatusError, "系统错误")

			context := log.System().NewContext(helper.GetTraceId(ctx))
			context["error"] = r
			context["trace"] = library.Trace()

			log.System().Error("system error", context)
		}
	}()
	ctx.Next()
}
