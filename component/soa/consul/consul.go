/**
集成consul组件，包含实例化consul客户端,服务发现,服务注册,服务注销,负载均衡等方法
*/
package consul

import (
	"fmt"
	"github.com/ebar-go/ego/component/soa/service"
	consulapi "github.com/hashicorp/consul/api"
)


var client *consulapi.Client

// InitClient 初始化consul客户端
func InitClient(config *consulapi.Config) error {
	var err error
	client, err =  consulapi.NewClient(config)

	return err
}

// DefaultConfig 默认配置
func DefaultConfig() *consulapi.Config {
	return consulapi.DefaultConfig()
}

// Discover 服务发现
func Discover(name string) (*service.Group, error) {
	services, _, err := client.Health().Service(name, "", true, &consulapi.QueryOptions{})
	if err != nil {
		return nil, fmt.Errorf("service: %s not found,%s", name, err.Error())
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("service name : %s not found", name)
	}

	group := new(service.Group)
	for _, item := range services {
		group.Add(service.Node{
			ID:      item.Service.ID,
			Name:    item.Service.Service,
			Address: item.Service.Address,
			Port:    item.Service.Port,
			Tags:    item.Service.Tags,
		})
	}

	return group, nil
}


// Register 注册服务
func Register(node service.Node) error {
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = node.ID
	registration.Name = node.Name
	registration.Port = node.Port
	registration.Tags = node.Tags
	registration.Address = node.Address

	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d%s", registration.Address, registration.Port, "/check")
	check.Timeout = "3s"
	check.Interval = "3s"
	check.DeregisterCriticalServiceAfter = "30s" //check失败后30秒删除本服务
	registration.Check = check

	return client.Agent().ServiceRegister(registration)
}

// Deregister 注销服务
func Deregister(node service.Node) error {
	return client.Agent().ServiceDeregister(node.ID)
}
