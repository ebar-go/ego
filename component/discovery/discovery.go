package discovery

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
