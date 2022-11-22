package etcd

import (
	"context"
	"fmt"
	"github.com/ebar-go/ego/utils/soa/etcd/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"log"
	"net"
	"testing"
	"time"
)

var (
	endpoints = []string{"127.0.0.1:2379"}
)

type Hello struct {
	pb.UnimplementedHelloServerServer
	addr string
}

func (h Hello) SayHello(ctx context.Context, in *pb.Req) (*pb.Resp, error) {
	return &pb.Resp{Name: "bar," + h.addr}, nil
}

func TestRegistry(t *testing.T) {
	addrs := []string{"127.0.0.1:8081", "127.0.0.1:8082", "127.0.0.1:8083"}
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

		registry := NewRegistry(Option{
			Endpoints:   endpoints,
			RegistryDir: "default",
			ServiceName: "app",
			ServiceAddr: addr,
			ServiceData: nil,
			Ttl:         time.Minute,
		})
		registry.Register()
	}

	select {}
}

func TestResolver(t *testing.T) {
	RegisterResolver("etcd", endpoints, "default", "app")
	cc, _ := grpc.DialContext(context.Background(), "etcd:///default/app", []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
	}...)
	client := pb.NewHelloServerClient(cc)
	resp, err := client.SayHello(context.Background(), &pb.Req{Name: "foo"})
	if err != nil {
		panic(err)
	}
	log.Println(resp.Name)
}
