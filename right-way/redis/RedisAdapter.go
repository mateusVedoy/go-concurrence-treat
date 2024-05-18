package redis

import (
	"fmt"
	"time"

	redis "github.com/go-redis/redis"
	redislock "github.com/stone-stones/redislock"
)

type RedisAdapter struct {
	client *redis.Client
	lock   *redislock.RedisLock
}

func (R *RedisAdapter) validate() error {
	_, err := R.client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func (R *RedisAdapter) Lock() error {
	key := "RIGHT_WAY_KEY"
	timeout := 15 * time.Second
	er := R.lock.Lock(key, timeout, false)
	if er != nil {
		return fmt.Errorf("Redis Adapter error locking key. Reason: %s", er.Error())
	}
	return nil
}

func (R *RedisAdapter) TryAccess() error {
	_, er := R.client.Get("KEY_TESTE").Result()

	if er != nil && er.Error() != "redis: nil" {
		return er
	}

	return nil
}

func (R *RedisAdapter) del() error {
	_, er := R.client.Del("KEY_TESTE").Result()

	if er != nil {
		return er
	}

	return nil
}

func (R *RedisAdapter) Unlock() error {
	er := R.del()
	if er != nil {
		return fmt.Errorf("Redis Adapter error unloking key. Reason: %s", er.Error())
	}
	return nil
}

func NewRedisAdapter() *RedisAdapter {

	adapter := RedisAdapter{}

	address := fmt.Sprintf("%s:%s", "127.0.0.1", "6379")

	redisCli := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})
	adapter.client = redisCli

	err := adapter.validate()

	if err != nil {
		panic(fmt.Errorf("Couldn't stablish redis connection. Reason: %s", err.Error()))
	}

	adapter.lock = redislock.NewRedisLock(adapter.client)

	return &adapter
}
