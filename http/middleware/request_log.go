package middleware

import (
	"bytes"
	"fmt"
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/ego/utils/date"
	"github.com/ebar-go/ego/utils/number"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
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
	requestTime := date.GetMicroTimeStampStr()
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	requestBody := getRequestBody(c)

	// 从头部信息获取
	traceId := strings.TrimSpace(c.GetHeader(app.Config().Server().TraceHeader))
	if traceId == "" {
		traceId = trace.NewId()
	}
	trace.SetTraceId(traceId)
	defer trace.DeleteTraceId()

	c.Next()

	// after request
	latency := time.Since(t)

	logContext := log.Context{}

	// 获取响应内容
	responseBody := blw.body.String()
	// 截断响应内容
	maxResponseSize := number.Min(number.Max(0, blw.body.Len()-1), app.Config().Server().MaxResponseLogSize)

	// 日志格式
	logContext["trace_id"] = traceId
	logContext["request_uri"] = c.Request.RequestURI
	logContext["request_method"] = c.Request.Method
	logContext["refer_service_name"] = c.Request.Referer()
	logContext["refer_request_host"] = c.ClientIP()
	logContext["request_body"] = requestBody
	logContext["request_time"] = requestTime
	logContext["response_time"] = date.GetMicroTimeStampStr()
	logContext["response_body"] = responseBody[0:maxResponseSize]
	logContext["time_used"] = fmt.Sprintf("%v", latency)
	logContext["header"] = c.Request.Header

	go app.LogManager().Request().Info("REQUEST LOG", logContext)
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
