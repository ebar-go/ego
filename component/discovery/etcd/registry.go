package etcd

import (
	"encoding/json"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"log"
	"time"
)

type Registry struct {
	etcd3Client *clientv3.Client
	key         string
	value       map[string]interface{}
	ttl         time.Duration
	ctx         context.Context
	cancel      context.CancelFunc
	leaseID     clientv3.LeaseID
}

type Option struct {
	RegistryDir string
	ServiceName string
	ServiceAddr string
	ServiceData map[string]interface{}
	Ttl         time.Duration
}

func NewRegistry(client *clientv3.Client, option Option) *Registry {
	ctx, cancel := context.WithCancel(context.Background())
	registry := &Registry{
		etcd3Client: client,
		key:         option.RegistryDir + "/" + option.ServiceName + "/" + option.ServiceAddr,
		value:       option.ServiceData,
		ttl:         option.Ttl / time.Second,
		ctx:         ctx,
		cancel:      cancel,
	}
	return registry
}

func (e *Registry) Register(stopCh <-chan struct{}) (err error) {
	if err = e.put(); err != nil {
		return
	}
	go func() {
		ticker := time.NewTicker(e.ttl / 5)
		for {
			select {
			case <-ticker.C:
				if lastErr := e.keepAlive(); lastErr != nil {
					log.Println("keepAlive failed", lastErr)
				}
			case <-stopCh:
				if lastErr := e.delete(); lastErr != nil {
					log.Println("delete failed", lastErr)
				}
				e.UnRegister()
				return
			}
		}
	}()

	return
}

func (e *Registry) UnRegister() {
	e.cancel()
	return
}

// =========================private methods =========================
func (e *Registry) put() error {
	resp, err := e.etcd3Client.Grant(e.ctx, int64(e.ttl))
	if err != nil {
		return errors.WithMessage(err, "grant key")
	}
	e.leaseID = resp.ID
	_, err = e.etcd3Client.Get(e.ctx, e.key)
	if err != nil && err != rpctypes.ErrKeyNotFound {
		return errors.WithMessage(err, "get key")
	}

	value, _ := json.Marshal(e.value)

	if err == nil || err == rpctypes.ErrKeyNotFound {
		_, err = e.etcd3Client.Put(e.ctx, e.key, string(value), clientv3.WithLease(resp.ID))
		return err
	}

	return errors.WithMessage(err, "get key")
}

func (e *Registry) delete() error {
	_, err := e.etcd3Client.Delete(e.ctx, e.key)
	return err
}

func (e *Registry) keepAlive() error {
	_, err := e.etcd3Client.KeepAliveOnce(e.ctx, e.leaseID)
	return err
}
