package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/examples/grpc/pb"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	_ "google.golang.org/grpc/encoding/gzip"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	version string
	// k8s环境下需要使用dns:///grpc-service:8081才能负载均衡
	target         string
	jaegerEndpoint string
)

func init() {
	flag.StringVar(&version, "version", "v1", "server version")
	flag.StringVar(&target, "target", "localhost:8081", "grpc server address")
	flag.StringVar(&jaegerEndpoint, "jaeger-endpoint", "http://jaeger-collector.istio-system.svc.cluster.local:14268/api/traces", "jaeger endpoint")
}

func main() {
	flag.Parse()
	aggregator := ego.New()

	aggregator.Aggregate(httpServer())
	aggregator.Aggregate(grpcServer())

	aggregator.Run()
}

func NewHttpTracer(serverName, address string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: serverName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{LogSpans: true, BufferFlushInterval: time.Second, CollectorEndpoint: address},
	}

	return cfg.NewTracer(config.Logger(jaegerlog.StdLogger))
}
func httpServer() runtime.Runnable {
	tracer, closer, err := NewHttpTracer("http-demo", jaegerEndpoint)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	cc, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithDefaultCallOptions(
			grpc.UseCompressor(gzip.Name),
		),
	)
	client := pb.NewUserServiceClient(cc)

	if err != nil {
		panic(err)
	}

	return ego.NewHTTPServer(":8080").
		EnableAvailableHealthCheck().
		EnablePprofHandler().
		WithDefaultRequestLogMiddleware().
		RegisterRouteLoader(func(router *gin.Engine) {
			router.Use(ginhttp.Middleware(tracer, ginhttp.OperationNameFunc(func(r *http.Request) string {
				return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.String())
			})))

			router.GET("greeter", func(ctx *gin.Context) {
				// new span
				span, traceCtx := opentracing.StartSpanFromContext(ctx.Request.Context(), "greeter")
				defer span.Finish()

				resp, err := client.Greet(traceCtx, &pb.GreetRequest{Name: "foo"})
				if err != nil {
					ctx.String(http.StatusInternalServerError, err.Error())
				} else {
					ctx.String(http.StatusOK, resp.Name)
				}

			})
		})
}

type UserService struct {
	pb.UnimplementedUserServiceServer
	tracer opentracing.Tracer
}

func (srv UserService) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	hostname, _ := os.Hostname()
	name := fmt.Sprintf("hostname=%s, version=%s", hostname, version)
	for i := 0; i < 10; i++ {
		name += name
	}
	return &pb.GreetResponse{Name: name}, nil
}

func grpcServer() runtime.Runnable {
	tracer, closer, err := NewHttpTracer("grpc-demo", jaegerEndpoint)
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	return ego.NewGRPCServer(":8081").RegisterService(func(s *grpc.Server) {
		pb.RegisterUserServiceServer(s, &UserService{tracer: tracer})
	})
}
