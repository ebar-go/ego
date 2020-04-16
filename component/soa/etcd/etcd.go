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
func InitClient(config client.Config) error   {
	var err error
	instance , err = client.New(config)
	return err
}

func Register(node service.Node) error {
	kapi := client.NewKeysAPI(instance)

	key := node.Name + "/" + node.ID
	val := net.JoinHostPort(node.Address, strconv.Itoa(node.Port))
	resp, err := kapi.Set(context.Background(), key, val, nil)

	fmt.Println(resp)
	return err
}

func Deregister(node service.Node) error {
	kapi := client.NewKeysAPI(instance)

	key := node.Name + "/" + node.ID
	resp, err := kapi.Delete(context.Background(), key, nil)

	fmt.Println(resp)
	return err
}

func Discover(name string) error {
	kapi := client.NewKeysAPIWithPrefix(instance, name)
	resp, err := kapi.Get(context.Background(), "*", nil)
	fmt.Println(resp)
	return err
}
