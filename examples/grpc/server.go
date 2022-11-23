package main

import (
	"context"
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/component"
	"github.com/ebar-go/ego/examples/grpc/pb"
	"github.com/ebar-go/ego/server"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func main() {
	aggregator := ego.New()

	aggregator.WithServer(httpServer(), grpcServer())

	aggregator.Run()
}

func httpServer() server.Server {
	return ego.NewHTTPServer(":8080").EnableAvailableHealthCheck().
		WithDefaultRequestLogMiddleware().RegisterRouteLoader(func(router *gin.Engine) {
		router.GET("greeter", func(ctx *gin.Context) {
			name := ctx.Query("name")
			cc, err := grpc.Dial("127.0.0.1:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				component.Provider().Logger().Errorf("NewClient: %v", err)
				return
			}
			resp, err := pb.NewUserServiceClient(cc).Greet(ctx, &pb.GreetRequest{Name: name})
			if err != nil {
				component.Provider().Logger().Errorf("Greet: %v", err)
				return
			}
			ctx.String(http.StatusOK, resp.Name)
		})
	})
}

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (srv UserService) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	return &pb.GreetResponse{Name: "greeter:" + in.GetName()}, nil
}

func grpcServer() server.Server {
	return ego.NewGRPCServer(":8081").RegisterService(func(s *grpc.Server) {
		pb.RegisterUserServiceServer(s, new(UserService))
	})
}
