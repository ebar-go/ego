package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ebar-go/ego/http/response"
	"fmt"
	"github.com/ebar-go/ego/http/constant"
	"github.com/ebar-go/ego/log"
	"github.com/ebar-go/ego/http/helper"
)

// NotFoundHandler 404
func NotFoundHandler(context *gin.Context)  {
	response.Error(context, constant.StatusNotFound, fmt.Sprintf("404 Not Found: %s", context.Request.RequestURI))
}

func Recover(ctx *gin.Context)  {
	defer func() {
		if r := recover(); r != nil {
			response.Error(ctx, constant.StatusError, "系统错误")
			log.System().Error("system error", log.System().NewContext(helper.GetTraceId(ctx)))
		}
	}()
	ctx.Next()
}