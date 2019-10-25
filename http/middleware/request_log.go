package middleware

import (
	"github.com/gin-gonic/gin"
	"bytes"
	"github.com/ebar-go/ego/library"
	"fmt"
	"time"
	"github.com/ebar-go/ego/log"
	"github.com/ebar-go/ego/http/request"
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

var accessChannel = make(chan string, 100)

// RequestLog gin的请求日志中间件
func RequestLog(c *gin.Context) {
	go handleAccessChannel()

	t := time.Now()
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	// 注册唯一ID
	traceId := request.GetTraceId(c)


	c.Next()


	// after request
	latency := time.Since(t)
	// 日志格式
	accessLogMap := make(map[string]interface{})

	accessLogMap["request_time"]      = latency
	accessLogMap["request_method"]    = c.Request.Method
	accessLogMap["request_uri"]       = c.Request.RequestURI
	accessLogMap["request_proto"]     = c.Request.Proto
	accessLogMap["request_ua"]        = c.Request.UserAgent()
	accessLogMap["request_referer"]   = c.Request.Referer()
	accessLogMap["request_post_data"] = c.Request.PostForm.Encode()
	accessLogMap["request_client_ip"] = c.ClientIP()
	accessLogMap["cost_time"] = fmt.Sprintf("%v", latency)
	accessLogMap["trace_id"] = traceId
	accessLogMap["response_body"] = blw.body.String()
	accessLogMap["status_code"] = blw.Status()

	accessLogJson, _ := library.JsonEncode(accessLogMap)
	fmt.Println(getArgs(c))

	accessChannel <- accessLogJson
}

// 这个函数只返回json化之后的数据，且不处理错误，错误就返回空字符串
func getArgs(c *gin.Context) string {
	if c.ContentType() == "multipart/form-data" {
		c.Request.ParseMultipartForm(32 << 20)
	} else {
		c.Request.ParseForm()
	}
	args, _ := library.JsonEncode(c.Request.Form)
	return args
}

func handleAccessChannel() {
	for accessLog := range accessChannel {
		log.GetSystemLogger().Info("request_log", accessLog)
	}

	return
}