package database

import (
	"boilerplate/utils"
	"boilerplate/utils/config"
	"boilerplate/utils/constant"
	cerror "boilerplate/utils/error"
	"boilerplate/utils/logger"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	ctx     context.Context
	client  *RedisClient
	clients map[string]*RedisClient
}

type RedisClient struct {
	prefix string
	client *redis.UniversalClient
}

func NewRedis(config *config.Config) *Redis {
	redis := &Redis{
		ctx:     config.Context,
		clients: map[string]*RedisClient{},
	}

	redis.Connect(config)

	return redis
}

func (r *Redis) Connect(config *config.Config) {
	for key, value := range config.Redis {
		client := redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:           value.Addrs,
			Username:        value.Username,
			Password:        value.Password,
			RouteByLatency:  value.Route.Latency,
			RouteRandomly:   value.Route.Random,
			DB:              value.DB,
			PoolSize:        value.PoolSize,
			PoolTimeout:     time.Duration(value.PoolTimeout) * time.Minute,
			MinIdleConns:    value.MinIdleConns,
			MaxIdleConns:    value.MaxIdleConns,
			ConnMaxIdleTime: time.Duration(value.ConnMaxIdleTime) * time.Minute,
			ConnMaxLifetime: time.Duration(value.ConnMaxLifeTime) * time.Minute,
		})

		r.clients[key] = &RedisClient{
			prefix: value.Prefix,
			client: &client,
		}

		if err := client.Ping(r.ctx).Err(); err != nil {
			logger.Sugar.Fatal(err)
		}
	}

	r.client = r.clients[constant.DEFAULT]

	logger.Sugar.Debug("Redis connected.")
}

func (r *Redis) Connection(name string) *Redis {
	if utils.IsEmptyString(name) {
		name = constant.DEFAULT
	}

	r.client = r.clients[name]

	return r
}

func (r *Redis) conn(key string) string {
	if r.client == nil {
		r.client = r.Connection(constant.DEFAULT).client
	}

	return r.prefixedKey(key)
}

func (r *Redis) prefixedKey(key string) string {
	prefix := r.client.prefix

	if !utils.IsEmptyString(prefix) {
		return fmt.Sprintf("%s:%s", prefix, key)
	}

	return key
}

func (r *Redis) Set(key string, data any, ttl time.Duration) (err error) {
	key = r.conn(key)

	if err = (*r.client.client).Set(r.ctx, key, data, ttl).Err(); err != nil {
		r.client = nil

		return cerror.Fail(cerror.FuncName(), "failed_redis_set", map[string]any{
			"redis_key":  key,
			"redis_data": data,
		}, err)
	}

	r.client = nil

	return nil
}

func (r *Redis) Get(key string) (result string, err error) {
	key = r.conn(key)

	if result, err = (*r.client.client).Get(r.ctx, key).Result(); err != nil {
		r.client = nil

		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", cerror.Fail(cerror.FuncName(), "failed_redis_get", map[string]any{"redis_key": key}, err)
	}

	r.client = nil

	return result, nil
}

func (r *Redis) Del(keys ...string) error {
	if r.client == nil {
		r.client = r.Connection(constant.DEFAULT).client
	}

	for idx, key := range keys {
		keys[idx] = r.prefixedKey(key)
	}

	if err := (*r.client.client).Del(r.ctx, keys...).Err(); err != nil {
		r.client = nil

		if errors.Is(err, redis.Nil) {
			return nil
		}

		return cerror.Fail(cerror.FuncName(), "failed_redis_remove", map[string]any{"redis_key": keys}, err)
	}

	r.client = nil

	return nil
}
