package etcd

import (
	"context"
	"github.com/ebar-go/ego/utils/soa/etcd/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
	"testing"
)

type T struct {
	pb.UnimplementedHelloServerServer
}

// 这部分就是proto文件编译后生成的签名格式，我们只需要按照自己的需要实现一下，是不是就像实现一个接口一样
func (t *T) SayHello(ctx context.Context, in *pb.Req) (*pb.Resp, error) {
	return &pb.Resp{Name: "hello" + in.Name}, nil

}

func TestNewService(t *testing.T) {
	addr := "127.0.0.1:8080"
	lis, err := net.Listen("tcp", addr)
	go func() {
		srv := grpc.NewServer()
		pb.RegisterHelloServerServer(srv, new(T))
		srv.Serve(lis)
	}()

	svc, err := NewService(ServiceInfo{Name: "app", IP: addr}, []string{"192.168.138.234:2379"})
	assert.Nil(t, err)
	log.Println(svc)

	err = svc.Start()
	assert.Nil(t, err)
	select {}

}

func TestResolver(t *testing.T) {
	resolver.Register(NewResolver([]string{"192.168.138.234:2379"}, "app"))
	client, err := grpc.DialContext(context.Background(), "etcd://app", []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}...)
	assert.Nil(t, err)

	cli := pb.NewHelloServerClient(client)
	resp, err := cli.SayHello(context.Background(), &pb.Req{Name: "hello"})
	assert.Nil(t, err)
	log.Println(resp)
}
