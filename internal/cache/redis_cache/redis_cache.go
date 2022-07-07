package redis_cache

import (
	"fmt"
	"github.com/abrbird/orders/config"
	"github.com/abrbird/orders/internal/cache"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	OrderPrefix = "order"
)

const (
	Nil = redis.Nil
)

type RedisCache struct {
	client *redis.Client
	order  *RedisOrderCache
}

func New(cfg config.Redis) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0, // use default DB
	})

	redisCache := &RedisCache{
		client: rdb,
	}

	redisCache.order = &RedisOrderCache{
		Prefix:     OrderPrefix,
		redisCache: redisCache,
		expiration: time.Minute * 15,
	}

	return redisCache
}

func (r RedisCache) Order() cache.OrderCache {
	return r.order
}
