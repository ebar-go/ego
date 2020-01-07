package consul

import (
	"fmt"
	"github.com/ebar-go/ego/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getClient() *Client {
	config := DefaultConfig()
	config.Address = "192.168.0.222:8500"

	return &Client{
		Config: config,
	}
}

func TestClient_Init(t *testing.T) {
	client := getClient()
	err := client.Init()
	assert.Nil(t, err)
	fmt.Println(client)
}

func TestClient_Register(t *testing.T) {
	client := getClient()

	ip, err := utils.GetLocalIp()
	assert.Nil(t, err)

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
	assert.Nil(t, err)
}

func TestClient_Discover(t *testing.T) {

	client := getClient()

	items, err := client.Discover("epet-go-demo")
	assert.Nil(t, err)
	fmt.Println(items)
}

func TestClient_LoadBalance(t *testing.T) {
	client := getClient()

	items, err := client.Discover("epet-go-demo")
	assert.Nil(t, err)

	service, err := client.LoadBalance(items)
	assert.Nil(t, err)
	fmt.Println(service.Address, service.ID)
}

func TestClient_DeRegister(t *testing.T) {
	t.SkipNow()

	client := getClient()

	err := client.DeRegister("epet-go-demo-2")
	assert.Nil(t, err)
}

func TestService_GetHost(t *testing.T) {

	client := getClient()

	items, err := client.Discover("epet-go-demo")
	assert.Nil(t, err)

	service, err := client.LoadBalance(items)
	fmt.Println(service.GetHost())
}
