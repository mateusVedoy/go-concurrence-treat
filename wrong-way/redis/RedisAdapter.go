package redis

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type RedisAdapter struct {
	client *redis.Client
}

func (R *RedisAdapter) Context() context.Context {
	return R.client.Context()
}

func (R *RedisAdapter) validate() error {
	_, err := R.client.Ping(R.Context()).Result()
	if err != nil {
		return err
	}
	return nil
}

func (R *RedisAdapter) Close() {
	R.client.Close()
}

func (R *RedisAdapter) Set(key string, timeout time.Duration) error {
	_, lockingErr := R.client.SetNX(R.Context(), key, "1", timeout).Result()

	if lockingErr != nil {
		return lockingErr
	}

	return nil
}

func (R *RedisAdapter) Del(key string) error {
	_, er := R.client.Del(R.Context(), key).Result()

	if er != nil {
		return er
	}

	return nil
}

func (R *RedisAdapter) Get(key string) (string, error) {
	result, er := R.client.Get(R.Context(), key).Result()

	if er != nil && er.Error() != "redis: nil" {
		return "", er
	}

	return result, nil
}

func NewRedisAdapter() *RedisAdapter {

	adapter := RedisAdapter{}

	address := fmt.Sprintf("%s:%s", "127.0.0.1", "6379")

	redisCli := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		Username: "",
		DB:       0,
	})
	adapter.client = redisCli

	err := adapter.validate()

	if err != nil {
		panic(fmt.Errorf("Couldn't stablish redis connection. Reason: %s", err.Error()))
	}

	return &adapter
}
