package soa

import (
	"fmt"
	"github.com/ebar-go/ego/utils/number"
	"github.com/ebar-go/ego/utils/soa/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

const (
	ETCDSchema = "etcd"
)

type ServiceInfo struct {
	Name   string
	Addr   string
	Weight int
}

type Instance interface {
	// Register a new service
	Register(stopCh <-chan struct{}, info ServiceInfo) error
	Discovery(serviceName string) (infos []ServiceInfo, err error)

	Resolver(serviceName string)
	BuildTarget(serviceName string) string
}

type ETCDDiscovery struct {
	endpoints []string
	namespace string
	ttl       time.Duration
	once      sync.Once
	cli       *clientv3.Client
}

func NewETCDDiscovery(endpoints []string, namespace string, ttl time.Duration) Instance {
	return &ETCDDiscovery{endpoints: endpoints, namespace: namespace, ttl: ttl}
}

func (discovery *ETCDDiscovery) getClient() *clientv3.Client {
	discovery.once.Do(func() {
		cli, err := clientv3.New(clientv3.Config{Endpoints: discovery.endpoints, DialTimeout: time.Second * 10})
		if err != nil {
			panic(err)
		}
		discovery.cli = cli

	})
	return discovery.cli
}

func (discovery *ETCDDiscovery) Register(stopCh <-chan struct{}, info ServiceInfo) error {
	return etcd.NewRegistry(discovery.getClient(), etcd.Option{
		RegistryDir: discovery.namespace,
		ServiceName: info.Name,
		ServiceAddr: info.Addr,
		ServiceData: map[string]interface{}{
			"weight": info.Weight,
		},
		Ttl: discovery.ttl,
	}).Register(stopCh)
}

func (discovery *ETCDDiscovery) Resolver(serviceName string) {
	etcd.RegisterResolver(ETCDSchema, discovery.getClient(), discovery.namespace, serviceName)
}

func (discovery *ETCDDiscovery) BuildTarget(serviceName string) string {
	return fmt.Sprintf("%s:///%s/%s", ETCDSchema, discovery.namespace, serviceName)
}

func (discovery *ETCDDiscovery) Discovery(serviceName string) (infos []ServiceInfo, err error) {
	items, err := etcd.Discovery(discovery.getClient(), discovery.namespace, serviceName)
	if err != nil {
		return
	}

	infos = make([]ServiceInfo, 0, len(items))
	for _, item := range items {
		weight := number.Int(item.Attributes.Value("weight"))
		infos = append(infos, ServiceInfo{Name: serviceName, Addr: item.Addr, Weight: weight})
	}
	return
}
