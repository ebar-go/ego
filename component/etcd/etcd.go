/**
 * @Author: Hongker
 * @Description:
 * @File:  etcd
 * @Version: 1.0.0
 * @Date: 2020/6/17 20:37
 */

package etcd

import (
	client "go.etcd.io/etcd/clientv3"
	"log"
)

// Client
type Client struct {
	instance *client.Client
	conf     *Config
}

func Connect(conf *Config) (*Client, error){
	instance, err := client.New(client.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: conf.Timeout,
	})
	if err != nil {
		return nil, err
	}

	log.Println("connect etcd success:", conf.Endpoints)

	return &Client{instance: instance}, nil
}
func (c *Client) Instance() *client.Client {
	return c.instance
}

func (c *Client) Api() client.KV {
	return client.NewKV(c.instance)
}
