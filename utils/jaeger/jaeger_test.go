package jaeger

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"io"
	"log"
	"net/http"
	"testing"
)

// start jaeger
//
//	docker run -d --name jaeger \
//	 -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
//	 -p 5775:5775/udp \
//	 -p 6831:6831/udp \
//	 -p 6832:6832/udp \
//	 -p 5778:5778 \
//	 -p 16686:16686 \
//	 -p 14268:14268 \
//	 -p 9411:9411 \
//	 jaegertracing/all-in-one:1.12
const (
	ginServerName  = "gin-server-demo"
	jaegerEndpoint = "127.0.0.1:6831"
)

func TestServer(t *testing.T) {
	tracer, closer, err := NewTracer(ginServerName, jaegerEndpoint)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	r := gin.Default()

	jaegerMiddle := ginhttp.Middleware(tracer, ginhttp.OperationNameFunc(func(r *http.Request) string {
		return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.String())
	}))
	r.Use(ginhttp.Middleware(tracer))
	r.Use(jaegerMiddle)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "pong",
		})
	})
	_ = r.Run(":8888")
}

const (
	// 服务名 服务唯一标示，服务指标聚合过滤依据。
	clientServerName = "demo-gin-client"
	ginEndpoint      = "http://127.0.0.1:8081"
)

func HandlerError(span opentracing.Span, err error) {
	span.SetTag(string(ext.Error), true)
	span.LogKV(otlog.Error(err))
	//log.Fatal("%v", err)
}
func TestClient(t *testing.T) {
	tracer, closer, err := NewTracer(clientServerName, jaegerEndpoint)
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	span := tracer.StartSpan("CallDemoServer")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	defer span.Finish()

	// 构建http请求
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/test", ginEndpoint),
		nil,
	)
	if err != nil {
		HandlerError(span, err)
		return
	}
	// 构建带tracer的请求
	req = req.WithContext(ctx)
	req, ht := nethttp.TraceRequest(tracer, req)
	defer ht.Finish()

	// 初始化http客户端
	httpClient := &http.Client{Transport: &nethttp.Transport{}}
	// 发起请求
	res, err := httpClient.Do(req)
	if err != nil {
		HandlerError(span, err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		HandlerError(span, err)
		return
	}
	log.Printf(" %s recevice: %s\n", clientServerName, string(body))
}
