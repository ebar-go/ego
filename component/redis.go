package component

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

type Redis struct {
	Named
	redis.UniversalClient
}

func (r *Redis) Connect(options *redis.Options) error {
	client := redis.NewClient(options)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	r.UniversalClient = client
	log.Println("Connect redis success:", options.Addr)
	return nil
}

func NewRedis() *Redis {
	return &Redis{Named: componentRedis}
}
