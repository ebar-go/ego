package soa

import (
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
