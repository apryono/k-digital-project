package redis

import "github.com/go-redis/redis"

type RedisClient struct {
	Client *redis.Client
}
