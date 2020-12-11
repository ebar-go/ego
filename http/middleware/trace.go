package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zutim/ego/app"
	"github.com/zutim/ego/component/event"
	"github.com/zutim/ego/component/trace"
	"strings"
)

func Trace(c *gin.Context) {
	// 从头部信息获取
	traceId := strings.TrimSpace(c.GetHeader(app.Config().Server().TraceHeader))
	if traceId == "" {
		traceId = trace.Id()
	}
	trace.Set(traceId)
	defer trace.GC()

	event.Trigger(event.BeforeRoute, nil)
	c.Next()
	event.Trigger(event.AfterRoute, nil)

}
