package etcd

import (
	"context"
	"github.com/ebar-go/ego/utils/soa/etcd/pb"
	"google.golang.org/grpc"
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
}

func (h Hello) SayHello(ctx context.Context, in *pb.Req) (*pb.Resp, error) {
	return &pb.Resp{Name: "bar"}, nil
}

func TestRegistry(t *testing.T) {
	addr := "127.0.0.1:8081"
	lis, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		panic(err)
	}

	go func() {
		srv := grpc.NewServer()
		pb.RegisterHelloServerServer(srv, new(Hello))
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
	select {}
}

func TestResolver(t *testing.T) {
	RegisterResolver("etcd", endpoints, "default", "app")
	cc, _ := grpc.DialContext(context.Background(), "etcd:///default/app", []grpc.DialOption{
		grpc.WithInsecure(),
	}...)
	client := pb.NewHelloServerClient(cc)
	resp, err := client.SayHello(context.Background(), &pb.Req{Name: "foo"})
	if err != nil {
		panic(err)
	}
	log.Println(resp.Name)
}
