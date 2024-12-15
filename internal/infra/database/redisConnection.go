package database

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var redisClienteIP *redis.Client
var redisClienteToken *redis.Client
var onceClienteIP sync.Once
var onceClienteToken sync.Once

func ObterRedisClienteIP() *redis.Client {
	onceClienteIP.Do(func() {
		redisClienteIP = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	})
	return redisClienteIP
}

func ObterRedisClienteToken() *redis.Client {
	onceClienteToken.Do(func() {
		redisClienteToken = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       1,
		})
	})
	return redisClienteToken
}
