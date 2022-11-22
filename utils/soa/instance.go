package soa

import (
	"context"
	"fmt"
	"github.com/ebar-go/ego/utils/soa/etcd"
	"time"
)

const (
	ETCDSchema = "etcd"
)

type Instance interface {
	// Register a new service
	Register(info ServiceInfo)
	Discovery(ctx context.Context, serviceName string) (infos []ServiceInfo, err error)

	Resolver(serviceName string)
	BuildTarget(serviceName string) string
}

type ETCDDiscovery struct {
	endpoints []string
	namespace string
	ttl       time.Duration
}

func NewETCDDiscovery(endpoints []string, namespace string, ttl time.Duration) Instance {
	return &ETCDDiscovery{endpoints: endpoints, namespace: namespace, ttl: ttl}
}

type ServiceInfo struct {
	Name   string
	Addr   string
	Weight int
}

func (discovery *ETCDDiscovery) Register(info ServiceInfo) {
	etcd.NewRegistry(etcd.Option{
		Endpoints:   discovery.endpoints,
		RegistryDir: discovery.namespace,
		ServiceName: info.Name,
		ServiceAddr: info.Addr,
		ServiceData: nil,
		Ttl:         discovery.ttl,
	}).Register()
}

func (discovery *ETCDDiscovery) Resolver(serviceName string) {
	etcd.RegisterResolver(ETCDSchema, discovery.endpoints, discovery.namespace, serviceName)
}

func (discovery *ETCDDiscovery) BuildTarget(serviceName string) string {
	return fmt.Sprintf("%s:///%s/%s", ETCDSchema, discovery.namespace, serviceName)
}

func (discovery *ETCDDiscovery) Discovery(ctx context.Context, serviceName string) (infos []ServiceInfo, err error) {
	items, err := etcd.Discovery(discovery.endpoints, discovery.namespace, serviceName)
	if err != nil {
		return
	}

	infos = make([]ServiceInfo, 0, len(items))
	for _, item := range items {
		weight, _ := item.Attributes.Value("weight").(int)
		infos = append(infos, ServiceInfo{Name: serviceName, Addr: item.Addr, Weight: weight})
	}
	return
}
