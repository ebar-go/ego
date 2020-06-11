package middleware

import (
	"github.com/ebar-go/ego/component/event"
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/ego/config"
	"github.com/gin-gonic/gin"
	"strings"
)

func Trace(c *gin.Context) {
	// 从头部信息获取
	traceId := strings.TrimSpace(c.GetHeader(config.Server().TraceHeader))
	if traceId == "" {
		traceId = trace.Id()
	}
	trace.Set(traceId)
	defer trace.GC()

	event.Trigger(event.BeforeRoute, nil)
	c.Next()
	event.Trigger(event.AfterRoute, nil)

}
