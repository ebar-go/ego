package soa

import (
	"context"
	"fmt"
	"github.com/ebar-go/ego/utils/runtime/signal"
	"github.com/ebar-go/ego/utils/soa/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"testing"
	"time"
)

var (
	endpoints = []string{"127.0.0.1:2379"}
	discovery = NewETCDDiscovery(endpoints, "default", time.Minute)
)

type Hello struct {
	pb.UnimplementedHelloServerServer
	addr string
}

func (h Hello) SayHello(ctx context.Context, in *pb.Req) (*pb.Resp, error) {
	return &pb.Resp{Name: "bar," + h.addr}, nil
}

func TestRegister(t *testing.T) {
	addrs := []string{"127.0.0.1:8081", "127.0.0.1:8082", "127.0.0.1:8083", "127.0.0.1:8084"}
	stopCh := signal.SetupSignalHandler()
	for _, addr := range addrs {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			panic(err)
		}

		go func() {
			srv := grpc.NewServer()
			pb.RegisterHelloServerServer(srv, &Hello{addr: addr})
			srv.Serve(lis)
		}()

		if err := discovery.Register(stopCh, ServiceInfo{Name: "app", Addr: addr}); err != nil {
			log.Println("register failed", err)
		}
	}

	<-stopCh
}

func TestResolver(t *testing.T) {
	discovery.Resolver("app")
	cc, _ := grpc.DialContext(context.Background(), discovery.BuildTarget("app"), []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithResolvers(),
	}...)
	client := pb.NewHelloServerClient(cc)
	for {
		resp, err := client.SayHello(context.Background(), &pb.Req{Name: "foo"})
		log.Println(resp, err)
		time.Sleep(time.Second)
	}

}
