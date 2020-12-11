package middleware

import (
	"bytes"
	"fmt"
	"github.com/ebar-go/ego/app"
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
func RequestLog(c *gin.Context) {
	t := time.Now()
	requestTime := egu.GetMicroTimeStampStr()
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	c.Next()

	// package log content
	items := log.Context{}
	items["request_uri"] = c.Request.RequestURI
	items["request_method"] = c.Request.Method
	items["refer_service_name"] = c.Request.Referer()
	items["refer_request_host"] = c.ClientIP()
	items["request_body"] = getRequestBody(c)
	items["request_time"] = requestTime
	items["response_time"] = egu.GetMicroTimeStampStr()
	items["response_body"] = getResponseBody(blw.body.String())
	items["time_used"] = fmt.Sprintf("%v", time.Since(t))

	if app.Config().Server().Debug {
		items["header"] = c.Request.Header
	}

	items["trace_id"] = trace.Get()

	// use goroutine
	go app.Logger().Info("REQUEST INFO", items)
}

// getResponseBody
func getResponseBody(s string) string {
	maxResponseSize := egu.Min(len(s), app.Config().Server().MaxResponseLogSize)
	res := make([]byte, maxResponseSize)
	copy(res, s[:maxResponseSize])
	return egu.Byte2Str(res)
}

// GetRequestBody 获取请求参数
func getRequestBody(c *gin.Context) interface{} {
	switch c.Request.Method {
	case http.MethodGet:
		return c.Request.URL.Query()
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodPatch:
		var bodyBytes []byte // 我们需要的body内容

		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return nil
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		return string(bodyBytes)
	}

	return nil
}
