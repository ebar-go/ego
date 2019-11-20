package middleware

import (
	"github.com/gin-gonic/gin"
	"bytes"
	"fmt"
	"time"
	"github.com/ebar-go/ego/log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/ego/http/constant"
	"strings"
	"github.com/ebar-go/ego/helper"
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
	requestTime := helper.GetTimeStampFloatStr()
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	requestBody := getRequestBody(c)

	// 从头部信息获取
	traceId := c.GetHeader(constant.GatewayTrace)
	if strings.TrimSpace(traceId) == "" {
		traceId = helper.NewTraceId()
	}
	trace.SetTraceId(traceId)
	defer trace.DeleteTraceId()

	c.Next()

	// after request
	latency := time.Since(t)

	logContext := log.Context{}


	responseBody := blw.body.String()
	maxRequestCount := helper.Min(blw.body.Len() - 1, constant.DefaultMaxResponseSize)
	// 日志格式
	logContext["trace_id"] = traceId
	logContext["request_uri"] = c.Request.RequestURI
	logContext["request_method"] = c.Request.Method
	logContext["refer_service_name"] = c.Request.Referer()
	logContext["refer_request_host"] = c.ClientIP()
	logContext["request_body"] = requestBody
	logContext["request_time"] = requestTime
	logContext["response_time"] = helper.GetTimeStampFloatStr()
	logContext["response_body"] = responseBody[0:maxRequestCount]
	logContext["time_used"] = fmt.Sprintf("%v", latency)
	logContext["header"] = c.Request.Header

	go log.Request().Info("REQUEST LOG", logContext)

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

		// 从原有Request.Body读取
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println(err)
		}

		// 新建缓冲区并替换原有Request.body
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		var params interface{}
		json.Unmarshal(bodyBytes, &params)
		return params

	}

	return nil
}

