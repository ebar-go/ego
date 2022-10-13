package jaeger

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"net/http"
	"time"
)

type Jaeger struct{}

type Tracer struct {
	instance opentracing.Tracer
	closer   io.Closer
}

func New(service, addr string) (*Tracer, error) {
	tracer, closer, err := NewOpenTracer(service, addr)
	if err != nil {
		return nil, err
	}

	return &Tracer{closer: closer, instance: tracer}, nil
}

func (tracer *Tracer) Instance() opentracing.Tracer {
	return tracer.instance
}

func (tracer *Tracer) Close() error {
	return tracer.closer.Close()
}

// NewSpan return opentracing.Span object, it should call Finish() method.
func (tracer *Tracer) NewSpan(name string) opentracing.Span {
	return tracer.instance.StartSpan("CallDemoServer")
}

func (tracer *Tracer) NewContext(ctx context.Context, span opentracing.Span) context.Context {
	return opentracing.ContextWithSpan(ctx, span)
}

func (tracer *Tracer) NewHttpRequestWithContext(ctx context.Context, req *http.Request) *nethttp.Tracer {
	var ht *nethttp.Tracer
	req, ht = nethttp.TraceRequest(tracer.instance, req.WithContext(ctx))
	return ht
}

func (tracer *Tracer) NewHttpRequestWithSpanName(name string, req *http.Request) *nethttp.Tracer {
	span := tracer.NewSpan(name)
	defer span.Finish()

	return tracer.NewHttpRequestWithContext(tracer.NewContext(context.Background(), span), req)
}

func (tracer *Tracer) ListenHttp(router *gin.Engine) {
	router.Use(ginhttp.Middleware(tracer.instance))
	router.Use(ginhttp.Middleware(tracer.instance, ginhttp.OperationNameFunc(func(r *http.Request) string {
		return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.String())
	})))
}

func NewOpenTracer(serverName, address string) (opentracing.Tracer, io.Closer, error) {
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
