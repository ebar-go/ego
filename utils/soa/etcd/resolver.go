package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"strings"
	"sync"
)

type etcdResolver struct {
	scheme        string
	etcdConfig    clientv3.Config
	etcdWatchPath string
	watcher       *Watcher
	cc            resolver.ClientConn
	wg            sync.WaitGroup
}

func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	etcdCli, err := clientv3.New(r.etcdConfig)
	if err != nil {
		return nil, err
	}
	r.cc = cc
	r.watcher = newWatcher(r.etcdWatchPath, etcdCli)
	r.start()
	return r, nil
}

func (r *etcdResolver) Scheme() string {
	return r.scheme
}

func (r *etcdResolver) start() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		out := r.watcher.Watch()
		for addrs := range out {
			if addrs != nil && len(addrs) > 0 {
				r.cc.UpdateState(resolver.State{Addresses: addrs})
			}
		}
	}()
}

func (r *etcdResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (r *etcdResolver) Close() {
	r.watcher.Close()
	r.wg.Wait()
}

func RegisterResolver(scheme string, endpoints []string, registryDir, srvName string) {
	resolver.Register(&etcdResolver{
		scheme:        scheme,
		etcdConfig:    clientv3.Config{Endpoints: endpoints},
		etcdWatchPath: parseRegistryDir(registryDir) + "/" + srvName,
	})
}

func parseRegistryDir(registryDir string) string {
	if registryDir == "" {
		return ""
	}
	if strings.Contains(registryDir, "/") {
		arr := strings.Split(registryDir, "/")
		for _, dir := range arr {
			if dir != "" {
				return dir
			}
		}
		return ""
	} else {
		return registryDir
	}
}
