package redis

import (
	"github.com/go-redis/redis"
	"log"
)

type Connection interface {
	GetInstance() redis.UniversalClient
}

type Single struct {
	conn redis.UniversalClient
}

func (s *Single) GetInstance() redis.UniversalClient{
	return s.conn
}

func Connect(conf *Config) (*Single, error) {
	connection := redis.NewClient(conf.Options())
	_, err := connection.Ping().Result()
	if err != nil {
		return nil, err
	}
	log.Println("Connect redis success:", conf.Host)

	return &Single{conn: connection}, nil
}