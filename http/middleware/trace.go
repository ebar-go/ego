package middleware

import (
	"github.com/ebar-go/ego/component/event"
	"github.com/ebar-go/ego/component/trace"
	"github.com/gin-gonic/gin"
	"strings"
)

func Trace(traceHeader string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从头部信息获取
		traceId := strings.TrimSpace(c.GetHeader(traceHeader))
		if traceId == "" {
			traceId = trace.Id()
		}
		trace.Set(traceId)
		defer trace.GC()

		event.Trigger(event.BeforeRoute, nil)
		c.Next()
		event.Trigger(event.AfterRoute, nil)
	}


}
