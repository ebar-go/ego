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
	"time"
)

// Client
type Client struct {
	instance *client.Client
	conf     *Config
}

func New(conf *Config) *Client {
	return &Client{conf: conf}
}

// Connect
func (c *Client) Connect() error {
	var err error
	c.instance, err = client.New(client.Config{
		Endpoints:   c.conf.Endpoints,
		DialTimeout: time.Second * time.Duration(c.conf.Timeout),
	})
	if err == nil {
		log.Println("Connect Etcd success:", c.conf.Endpoints)
	}

	return err
}

func (c *Client) Instance() *client.Client {
	return c.instance
}

func (c *Client) Api() client.KV {
	return client.NewKV(c.instance)
}
