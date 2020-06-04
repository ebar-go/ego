package consul

import (
	"fmt"
	"github.com/ebar-go/ego/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(m *testing.M) {
	config := DefaultConfig()
	config.Address = "10.0.75.2:8500"
	err := InitClient(config)
	fmt.Println(err)
	m.Run()
}

func TestClient_Register(t *testing.T) {
	ip, err := utils.GetLocalIp()
	assert.Nil(t, err)

	registration := NewServiceRegistration()
	registration.ID = "service-id"
	registration.Name = "serviceName"
	registration.Port = 8085
	registration.Tags = []string{"test service"}
	registration.Address = ip

	check := NewServiceCheck()
	check.HTTP = fmt.Sprintf("http://%s:%d%s", registration.Address, registration.Port, "/check")
	check.Timeout = "3s"
	check.Interval = "3s"
	check.DeregisterCriticalServiceAfter = "30s" //check失败后30秒删除本服务
	registration.Check = check

	err = Register(registration)
	fmt.Println(err)
	assert.Nil(t, err)
}

func TestClient_Discover(t *testing.T) {

	group, err := Discover("serviceName")

	assert.Nil(t, err)

	fmt.Println(group.First())

}

func TestClient_Deregister(t *testing.T) {
	t.SkipNow()

	err := Deregister("serviceName")
	assert.Nil(t, err)
}
