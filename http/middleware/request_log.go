package middleware

import (
	"github.com/gin-gonic/gin"
	"bytes"
	"github.com/ebar-go/ego/library"
	"fmt"
	"time"
	"github.com/ebar-go/ego/log"
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
	accessLogMap["response_body"] = blw.body.String()
	accessLogMap["status_code"] = blw.Status()

	accessLogJson, _ := library.JsonEncode(accessLogMap)
	accessChannel <- accessLogJson
}

func handleAccessChannel() {
	for accessLog := range accessChannel {
		log.RequestLogger.Info("request_log", accessLog)
	}

	return
}