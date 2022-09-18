package component

const (
	defaultWeight = 100
)

// Discovery manage service register and discover.
type Discovery interface {
	Discover(name string) (*Service, error)
	Register(service Service)
	Unregister(service Service)
}

type Service struct {
	name   string
	weight int
}

func (s *Service) WithWeight(weight int) *Service {
	s.weight = weight
	return s
}

func NewService(name string) *Service {
	return &Service{name: name, weight: defaultWeight}
}

type EtcdDiscovery struct{}
type FileDiscovery struct{}
