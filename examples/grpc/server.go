package main

import (
	"context"
	"fmt"
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/examples/grpc/pb"
	"github.com/ebar-go/ego/utils/runtime"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"strconv"
)

func main() {
	aggregator := ego.New()

	aggregator.Aggregate(httpServer())
	aggregator.Aggregate(grpcServer())

	aggregator.Run()
}

const (
	target = "localhost:8081" // 与本地grpc服务建立连接，负载均衡无效
	//target = "grpc-service:8081" // 与k8s的service建立连接，负载均衡无效
	//target = "dns:///grpc-service:8081" // 与k8s的headless service建立连接，负载均衡有效
)

func httpServer() runtime.Runnable {
	cc, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
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
			router.GET("greeter", func(ctx *gin.Context) {

				num, _ := strconv.Atoi(ctx.Query("num"))
				if num == 0 {
					num = 10
				}

				for i := 0; i < num; i++ {
					resp, err := client.Greet(ctx, &pb.GreetRequest{Name: "foo"})
					if err != nil {
						component.Logger().Errorf("Greet: %v", err)
						return
					}
					component.Logger().Info(resp.Name)
				}

				//ctx.String(http.StatusOK, resp.Name)
			})
		})
}

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (srv UserService) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	hostname, _ := os.Hostname()
	return &pb.GreetResponse{Name: fmt.Sprintf("hostname=%s, in=%s, out=%s", hostname, in.Name, "bar")}, nil
}

func grpcServer() runtime.Runnable {
	return ego.NewGRPCServer(":8081").RegisterService(func(s *grpc.Server) {
		pb.RegisterUserServiceServer(s, new(UserService))
	})
}
