package redis

import (
	"errors"
	"fmt"
	"time"
)

type RedisLockingService struct {
	redisAdapter    *RedisAdapter
	LOCKING_KEY     string
	TIMEOUT_SECONDS time.Duration
}

func (R *RedisLockingService) Lock() error {
	return R.redisAdapter.Set(R.LOCKING_KEY, R.TIMEOUT_SECONDS)
}

func (R *RedisLockingService) Unlock() error {
	return R.redisAdapter.Del(R.LOCKING_KEY)
}

func (R *RedisLockingService) TryAccess() error {
	data, er := R.redisAdapter.Get(R.LOCKING_KEY)

	if er != nil || data != "" {
		return fmt.Errorf("access denied by resource locked, try again in a few seconds")
	}

	return nil
}

func NewRedisLockingService(redisAdapter *RedisAdapter) *RedisLockingService {

	key := "WRONG_WAY_KEY"

	if key == "" {
		panic(errors.New("REDIS_LOCKING_KEY must be informed"))
	}

	timeout := 15 * time.Second

	return &RedisLockingService{
		redisAdapter:    redisAdapter,
		LOCKING_KEY:     key,
		TIMEOUT_SECONDS: timeout,
	}
}
