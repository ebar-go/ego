package jaeger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"net/http"
	"time"
)

type Jaeger struct{}

func NewTracer(serverName, address string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: serverName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{LogSpans: true, BufferFlushInterval: time.Second},
	}

	transport, err := jaeger.NewUDPTransport(address, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(transport)
	return cfg.NewTracer(config.Reporter(reporter))
}

func WithGinEngine(router *gin.Engine, tracer opentracing.Tracer) {
	router.Use(ginhttp.Middleware(tracer))
	router.Use(ginhttp.Middleware(tracer, ginhttp.OperationNameFunc(func(r *http.Request) string {
		return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.String())
	})))
}
