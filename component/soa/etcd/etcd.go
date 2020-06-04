package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/ebar-go/ego/component/soa/service"
	"net"
	"strconv"
)

var instance client.Client

// InitClient 初始化Client
func InitClient(config client.Config) error {
	var err error
	instance, err = client.New(config)
	return err
}

// Register 服务注册
func Register(node service.Node) error {
	kapi := client.NewKeysAPI(instance)

	key := node.Name + "/" + node.ID
	val := net.JoinHostPort(node.Address, strconv.Itoa(node.Port))
	resp, err := kapi.Set(context.Background(), key, val, nil)

	fmt.Println(resp)
	return err
}

// Deregister 服务注销
func Deregister(node service.Node) error {
	kapi := client.NewKeysAPI(instance)

	key := node.Name + "/" + node.ID
	resp, err := kapi.Delete(context.Background(), key, nil)

	fmt.Println(resp)
	return err
}

// Discover 服务发现
func Discover(name string) error {
	kapi := client.NewKeysAPIWithPrefix(instance, name)
	resp, err := kapi.Get(context.Background(), "*", nil)
	fmt.Println(resp)
	return err
}
