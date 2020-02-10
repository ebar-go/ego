/**
集成consul组件，包含实例化consul客户端,服务发现,服务注册,服务注销,负载均衡等方法
*/
package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// Service 服务
type Service struct {
	// 服务ID
	ID string

	// 服务名称
	Name string

	// 服务地址
	Address string

	// 服务端口
	Port int

	// 服务标签
	Tags []string

	Weight int
}

func (server *Service) GetWeight() int {
	return server.Weight
}

// Client 客户端
type Client struct {
	Config       *consulapi.Config // 配置
	consulClient *consulapi.Client // consul客户端
	initialize   bool
}

// DefaultConfig 默认配置
func DefaultConfig() *consulapi.Config {
	return consulapi.DefaultConfig()
}

// NewServiceRegistration 实例化服务注册项
func NewServiceRegistration() *consulapi.AgentServiceRegistration {
	r := new(consulapi.AgentServiceRegistration)
	r.Weights = &consulapi.AgentWeights{Passing:1, Warning:1}
	return r
}

// NewServiceCheck 实例化服务检查项
func NewServiceCheck() *consulapi.AgentServiceCheck {
	return new(consulapi.AgentServiceCheck)
}

// NewClient 获取客户端
func (client *Client) Init() (err error) {
	if client.initialize {
		return fmt.Errorf("请勿重复初始化")
	}

	if client.Config == nil {
		return fmt.Errorf("缺少配置项")
	}

	client.consulClient, err = consulapi.NewClient(client.Config)
	if err != nil {
		return fmt.Errorf("初始化失败:%s", err.Error())
	}

	client.initialize = true
	return nil
}

// instance 获取实例
func (client *Client) instance() *consulapi.Client {
	if !client.initialize {
		err := client.Init()
		if err != nil {
			panic(err)
		}
	}

	return client.consulClient
}

// Discover 服务发现
func (client *Client) Discover(name string) ([]Service, error) {
	services, _, err := client.instance().Health().Service(name, "", true, &consulapi.QueryOptions{})
	if err != nil {
		return nil, fmt.Errorf("service name : %s not found,%s", name, err.Error())
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("service name : %s not found", name)
	}

	var serviceItems []Service
	for _, service := range services {
		serviceItem := Service{
			ID:      service.Service.ID,
			Name:    service.Service.Service,
			Address: service.Service.Address,
			Port:    service.Service.Port,
			Tags:    service.Service.Tags,
			Weight: service.Service.Weights.Warning,
		}
		serviceItems = append(serviceItems, serviceItem)
	}

	return serviceItems, nil
}

// LoadBalance 负载均衡
func (client *Client) LoadBalance(serviceItems []Service) (*Service, error) {
	if len(serviceItems) == 0 {
		return nil, fmt.Errorf("Serivce items is empty")
	}

	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator

	return &serviceItems[rand.Intn(len(serviceItems))], nil
}

// Register 注册
func (client *Client) Register(registration *consulapi.AgentServiceRegistration) error {
	return client.instance().Agent().ServiceRegister(registration)
}

// DeRegister 注销服务
func (client *Client) DeRegister(name string) error {
	return client.instance().Agent().ServiceDeregister(name)
}

// GetHost 获取服务完整的地址
func (service *Service) GetHost() string {
	return net.JoinHostPort(service.Address, strconv.Itoa(service.Port))
}
