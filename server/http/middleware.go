package http

import (
	"fmt"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// CORS 跨域中间件
func CORS(c *gin.Context) {
	method := c.Request.Method

	// set response header
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

func Trace(traceHeader string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从头部信息获取
		traceId := strings.TrimSpace(c.GetHeader(traceHeader))
		if traceId != "" {
			component.Tracer().Set(traceId)
		}
		defer component.Tracer().Release()

		c.Next()
	}

}

func RequestLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()
		ctx.Next()

		items := gin.H{}
		items["request_uri"] = ctx.Request.RequestURI
		items["request_method"] = ctx.Request.Method
		items["refer_service_name"] = ctx.Request.Referer()
		items["refer_request_host"] = ctx.ClientIP()
		items["time_used"] = fmt.Sprintf("%v", time.Since(t))
		component.Logger().Infof("request log: %v", items)
	}
}

func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer runtime.HandleCrash()
		ctx.Next()
	}
}
