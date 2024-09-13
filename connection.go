package cacheproxy

import (
	"context"
	"fmt"
	"log"

	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	ctx context.Context
	c   *redis.Client
}

const CACHE_TTL = 120 * time.Second

var ctx context.Context
var cancelCtx context.CancelFunc

func initRedisConnection(dns string, ctx context.Context) (*RedisClient, error) {
	c := redis.NewClient(&redis.Options{
		Addr:             dns,
		Password:         "",
		DB:               6,
		DisableIndentity: true,
	})

	connectionError := checkRedisAvailability(c, ctx)
	if connectionError != nil {
		return nil, connectionError
	}

	return &RedisClient{
		ctx: ctx,
		c:   c,
	}, nil
}

func checkRedisAvailability(c *redis.Client, ctx context.Context) error {
	_, err := c.Ping(ctx).Result()
	if err == nil {
		return nil
	}

	e := fmt.Sprintf("redis instannce connection error: %v", err)
	return errors.New(e)
}

func (p *Plugin) Weight() uint {
	return 100
}

func (r *RedisClient) Get(key string) ([]byte, error) {
	val, err := r.c.Get(r.ctx, key).Bytes()

	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		// Handle other Redis errors
		log.Printf("Error fetching key %s: %v", key, err)
		return nil, err
	}

	return val, nil
}

func (r *RedisClient) Set(key string, v interface{}) {
	r.c.Set(r.ctx, key, v, CACHE_TTL)
}
