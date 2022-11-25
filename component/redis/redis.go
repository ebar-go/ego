package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type Instance struct {
	redis.UniversalClient
}

func (r *Instance) Connect(options *redis.Options) error {
	client := redis.NewClient(options)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	r.UniversalClient = client
	return nil
}

func New() *Instance {
	return &Instance{}
}
