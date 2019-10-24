package consul

import (
	"testing"
	"git.epetbar.com/go-package/ego/test"
	"fmt"
	"git.epetbar.com/go-package/ego/library"
)

func getClient() *Client {
	config := DefaultConfig()
	config.Address = "192.168.0.222:8500"

	return &Client{
		Config:config,
	}
}

func TestClient_Init(t *testing.T) {
	client := getClient()
	err := client.Init()
	test.AssertNil(t, err)
	fmt.Println(client)
}

func TestClient_Register(t *testing.T) {
	client := getClient()

	ip, err := library.GetLocalIp()
	test.AssertNil(t, err)

	registration := NewServiceRegistration()
	registration.ID = "epet-go-demo-2"
	registration.Name = "epet-go-demo"
	registration.Port = 8088
	registration.Tags = []string{"epet-go-demo"}
	registration.Address = ip

	check := NewServiceCheck()
	check.HTTP = fmt.Sprintf("http://%s:%d%s", registration.Address, registration.Port, "/check")
	check.Timeout = "3s"
	check.Interval = "3s"
	check.DeregisterCriticalServiceAfter = "30s" //check失败后30秒删除本服务
	registration.Check = check

	err = client.Register(registration)
	fmt.Println(err)
	test.AssertNil(t, err)
}

func TestClient_Discover(t *testing.T) {

	client := getClient()

	items, err := client.Discover("epet-go-demo")
	test.AssertNil(t, err)
	fmt.Println(items)
}

func TestClient_LoadBalance(t *testing.T) {
	client := getClient()

	items, err := client.Discover("epet-go-demo")
	test.AssertNil(t, err)

	service, err := client.LoadBalance(items)
	test.AssertNil(t, err)
	fmt.Println(service.Address, service.ID)
}

func TestClient_DeRegister(t *testing.T) {
	t.SkipNow()

	client := getClient()


	err := client.DeRegister("epet-go-demo-2")
	test.AssertNil(t, err)
}

func TestService_GetHost(t *testing.T) {

	client := getClient()

	items, err := client.Discover("epet-go-demo")
	test.AssertNil(t, err)

	service, err := client.LoadBalance(items)
	fmt.Println(service.GetHost())
}