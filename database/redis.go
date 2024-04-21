package database

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	client  *redis.Client
	context context.Context
}

func (r *RedisClient) NewClient(clientConfiguration *ClientConfiguration) DatabaseClient {
	r.client = redis.NewClient(&redis.Options{
		Addr:     clientConfiguration.Address,
		Password: clientConfiguration.Password,
		DB:       0,
	})
	r.context = clientConfiguration.Context
	return r
}

func (r *RedisClient) HealthCheck() error {
	return r.client.Ping().Err()
}

func (r *RedisClient) Get(key string) (interface{}, error) {
	return r.client.Get(key).Result()
}
func (r *RedisClient) Set(key string, value interface{}) error {
	return r.client.Set(key, value, (24*7)*time.Hour).Err()
}
