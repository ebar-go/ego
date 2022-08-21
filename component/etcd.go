package component

import (
	client "go.etcd.io/etcd/clientv3"
	"time"
)

type Etcd struct {
	Named

	instance *client.Client
}

func (etcd *Etcd) Connect(endpoints []string, timeout time.Duration) ( error) {
	instance, err := client.New(client.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})

	etcd.instance = instance

	return err
}

func NewEtcd() *Etcd {
	return &Etcd{Named: "etcd"}
}
