package redis

import (
	"github.com/go-redis/redis"
	"log"
)

type Connection interface {
	GetInstance() redis.UniversalClient
}

type conn struct {
	client redis.UniversalClient
}

func (c *conn) GetInstance() redis.UniversalClient {
	return c.client
}

func Connect(conf *Config) (Connection, error) {
	client := redis.NewClient(conf.Options())
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	log.Println("Connect redis success:", conf.Host)

	return &conn{client: client}, nil
}
