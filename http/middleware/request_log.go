package middleware

import (
	"github.com/gin-gonic/gin"
	"bytes"
	"github.com/ebar-go/ego/library"
	"fmt"
	"time"
	"github.com/ebar-go/ego/log"
	"github.com/ebar-go/ego/http/helper"
)

// bodyLogWriter 读取响应Writer
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 读取响应数据
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

var accessChannel = make(chan log.Context, 100)

// RequestLog gin的请求日志中间件
func RequestLog(c *gin.Context) {
	go handleAccessChannel()

	t := time.Now()
	requestTime := library.GetTimeStampFloatStr()
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	requestBody := helper.GetRequestBody(c)
	c.Next()

	// after request
	latency := time.Since(t)

	logContext := log.NewContext(helper.GetTraceId(c))
	// 日志格式
	logContext["request_uri"] = c.Request.RequestURI
	logContext["request_method"] = c.Request.Method
	logContext["refer_service_name"] = c.Request.Referer()
	logContext["refer_request_host"] = c.ClientIP()
	logContext["request_body"] = requestBody
	logContext["request_time"] = requestTime
	logContext["response_time"] = library.GetTimeStampFloatStr()
	logContext["response_body"] = blw.body.String()
	logContext["time_used"] = fmt.Sprintf("%v", latency)

	logContext["header"] = c.Request.Header

	accessChannel <- logContext
}

func handleAccessChannel() {
	for accessLog := range accessChannel {

		log.Request().Info("REQUEST LOG", accessLog)
	}

	return
}
