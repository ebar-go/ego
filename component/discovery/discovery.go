package discovery

import "sync"

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

var discoveryInstance = struct {
	once     sync.Once
	instance Instance
}{}

func SetInstance(instance Instance) {
	discoveryInstance.once.Do(func() {
		discoveryInstance.instance = instance
	})
}

func Get() Instance {
	return discoveryInstance.instance
}
