package middleware

import (
	"bytes"
	"fmt"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/egu"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
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

// RequestLog gin的请求日志中间件
func RequestLog(logger *log.Logger, maxResponseSize int) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		t := time.Now()
		requestTime := egu.GetMicroTimeStampStr()
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		ctx.Next()

		// package log content
		items := log.Context{}
		items["request_uri"] = ctx.Request.RequestURI
		items["request_method"] = ctx.Request.Method
		items["refer_service_name"] = ctx.Request.Referer()
		items["refer_request_host"] = ctx.ClientIP()
		items["request_body"] = getRequestBody(ctx)
		items["request_time"] = requestTime
		items["response_time"] = egu.GetMicroTimeStampStr()
		items["response_body"] = getResponseBody(blw.body.String(), maxResponseSize)
		items["time_used"] = fmt.Sprintf("%v", time.Since(t))

		// use goroutine
		trace.Go(func() {
			logger.Info("request_log", items)
		})
	}
	
}

// getResponseBody
func getResponseBody(s string, maxLogSize int) string {
	maxResponseSize := egu.Min(len(s), maxLogSize)
	res := make([]byte, maxResponseSize)
	copy(res, s[:maxResponseSize])
	return egu.Byte2Str(res)
}

// GetRequestBody 获取请求参数
func getRequestBody(ctx *gin.Context) interface{} {
	switch ctx.Request.Method {
	case http.MethodGet:
		return ctx.Request.URL.Query()
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodPatch:
		var bodyBytes []byte // 我们需要的body内容

		bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			return nil
		}
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		return string(bodyBytes)
	}

	return nil
}
