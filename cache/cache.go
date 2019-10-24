package cache

import (
	"errors"
	"github.com/go-redis/redis/v7"
	"time"
)

var KeyNotExisted = errors.New("key is not existed")

type Cache interface {
	Set(key string, value string, expiredTime time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
}

type redisCache struct {
	addr string
	rdb  *redis.Client
}

func New(addr string) Cache {
	options := &redis.Options{
		Addr:         addr,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	}

	return &redisCache{
		addr: addr,
		rdb:  redis.NewClient(options),
	}
}

func (c *redisCache) Set(key string, value string, expire time.Duration) error {
	return c.rdb.Set(key, value, expire).Err()
}

func (c *redisCache) Get(key string) (string, error) {
	val, err := c.rdb.Get(key).Result()
	if err == redis.Nil {
		return "", KeyNotExisted
	}

	if err != nil {
		return "", err
	}

	return val, nil
}

func (c *redisCache) Del(key string) error {
	_, err := c.rdb.Del(key).Result()
	return err
}
