package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	"strings"
	"sync"
)

type Watcher struct {
	key    string
	client *clientv3.Client
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	addrs  []resolver.Address
}

func (w *Watcher) Close() {
	w.cancel()
}

func newWatcher(key string, cli *clientv3.Client) *Watcher {
	ctx, cancel := context.WithCancel(context.Background())
	w := &Watcher{
		key:    key,
		client: cli,
		ctx:    ctx,
		cancel: cancel,
	}
	return w
}

func (w *Watcher) GetAllAddresses() []resolver.Address {
	ret := make([]resolver.Address, 0)

	resp, err := w.client.Get(w.ctx, w.key, clientv3.WithPrefix())
	if err == nil {
		return extractAddresses(resp)
	}
	return ret
}

func (w *Watcher) Watch() chan []resolver.Address {
	out := make(chan []resolver.Address, 10)
	w.wg.Add(1)
	go func() {
		defer func() {
			close(out)
			w.wg.Done()
		}()
		w.addrs = w.GetAllAddresses()
		out <- w.cloneAddresses(w.addrs)

		rch := w.client.Watch(w.ctx, w.key, clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				switch ev.Type {
				case mvccpb.PUT:
					address := resolver.Address{}
					if err := convertEventToAddress(ev.Kv, &address); err != nil {
						grpclog.Infoln(err)
						continue
					}
					if w.addAddr(address) {
						out <- w.cloneAddresses(w.addrs)
					}
				case mvccpb.DELETE:
					address := resolver.Address{}
					if err := convertEventToAddress(ev.Kv, &address); err != nil {
						grpclog.Infoln(err)
						continue
					}
					if w.removeAddr(address) {
						out <- w.cloneAddresses(w.addrs)
					}
				}
			}
		}
	}()
	return out
}

func convertEventToAddress(kv *mvccpb.KeyValue, address *resolver.Address) error {
	//获取地址("ip:port")
	arr := strings.Split(string(kv.Key), "/")
	if arr == nil || len(arr) == 0 {
		return errors.New("invalid key")
	}
	last := arr[len(arr)-1]
	lastArr := strings.Split(last, "-")
	if lastArr == nil || len(lastArr) == 0 {
		return errors.New("invalid key")
	}
	address.Addr = lastArr[len(lastArr)-1]
	address.Attributes = attributes.New("someKey", "someValue")
	//获取属性
	attr := make(map[string]interface{}, 0)
	if err := json.Unmarshal(kv.Value, &attr); err != nil && len(attr) > 0 {
		for attrKey, attrValue := range attr {
			address.Attributes.WithValue(attrKey, attrValue)
		}
	}
	return nil
}

func extractAddresses(resp *clientv3.GetResponse) []resolver.Address {
	addresses := make([]resolver.Address, 0)
	if resp == nil || resp.Kvs == nil {
		return addresses
	}
	for i := range resp.Kvs {
		address := resolver.Address{}
		if k := resp.Kvs[i].Key; k != nil {
			arr := strings.Split(string(k), "/")
			address.Addr = arr[len(arr)-1]
		} else {
			continue
		}
		address.Attributes = attributes.New("someKey", "someValue")
		if v := resp.Kvs[i].Value; v != nil {
			attr := make(map[string]interface{}, 0)
			if err := json.Unmarshal(v, &attr); err != nil && len(attr) > 0 {
				for attrKey, attrValue := range attr {
					address.Attributes.WithValue(attrKey, attrValue)
				}
			}
		}
		addresses = append(addresses, address)
	}
	return addresses
}

func (w *Watcher) cloneAddresses(in []resolver.Address) []resolver.Address {
	out := make([]resolver.Address, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = in[i]
	}
	return out
}

func (w *Watcher) addAddr(addr resolver.Address) bool {
	for _, v := range w.addrs {
		if addr.Addr == v.Addr {
			return false
		}
	}
	w.addrs = append(w.addrs, addr)
	return true
}

func (w *Watcher) removeAddr(addr resolver.Address) bool {
	for i, v := range w.addrs {
		if addr.Addr == v.Addr {
			w.addrs = append(w.addrs[:i], w.addrs[i+1:]...)
			return true
		}
	}
	return false
}
