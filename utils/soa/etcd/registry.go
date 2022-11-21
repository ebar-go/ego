package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

// 服务信息
type ServiceInfo struct {
	Name string
	IP   string
}

type Service struct {
	ServiceInfo ServiceInfo
	stop        chan error
	leaseId     clientv3.LeaseID
	client      *clientv3.Client
}

// NewService 创建一个注册服务
func NewService(info ServiceInfo, endpoints []string) (service *Service, err error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 200,
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	service = &Service{
		ServiceInfo: info,
		client:      client,
	}
	return
}

// Start 注册服务启动
func (service *Service) Start() (err error) {

	ch, err := service.keepAlive()
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		select {
		case err := <-service.stop:
			return err
		case <-service.client.Ctx().Done():
			return errors.New("service closed")
		case resp, ok := <-ch:
			// 监听租约
			if !ok {
				log.Println("keep alive channel closed")
				return service.revoke()
			}
			log.Printf("Recv reply from service: %s, ttl:%d", service.getKey(), resp.TTL)
		}
	}

	return
}

func (service *Service) Stop() {
	service.stop <- nil
}

func (service *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	info := &service.ServiceInfo
	key := info.Name + "/" + info.IP
	val, _ := json.Marshal(info)

	// 创建一个租约
	resp, err := service.client.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = service.client.Put(context.TODO(), key, string(val), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	service.leaseId = resp.ID
	return service.client.KeepAlive(context.TODO(), resp.ID)
}

func (service *Service) revoke() error {
	_, err := service.client.Revoke(context.TODO(), service.leaseId)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("servide:%s stop\n", service.getKey())
	return err
}

func (service *Service) getKey() string {
	return service.ServiceInfo.Name + "/" + service.ServiceInfo.IP
}
