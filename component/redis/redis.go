package redis

import (
	goredis "github.com/go-redis/redis"
	"log"
)

// Client
type Client struct {
	conf *Config
	goredis.UniversalClient
}

func New(conf *Config) *Client {
	return &Client{conf:conf}
}


// Connect 单点连接
func (client *Client) Connect() error {
	connection := goredis.NewClient(client.conf.Options())
	_, err := connection.Ping().Result()
	if err != nil {
		return err
	}
	log.Printf("redis connect success:%s", client.conf.Host)
	client.UniversalClient = connection
	return nil
}

// ConnectCluster 连接集群
func (client *Client) ConnectCluster() error {
	connection := goredis.NewClusterClient(client.conf.ClusterOption())
	_, err := connection.Ping().Result()
	if err != nil {
		return err
	}
	log.Printf("redis connect success:%s", client.conf.Cluster)
	client.UniversalClient = connection

	return nil
}
