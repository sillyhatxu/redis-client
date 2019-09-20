package redislock

import (
	"fmt"
	"github.com/sillyhatxu/redis-client/redis"
	"sync"
	"time"
)

const (
	Lock   = "1"
	Unlock = "0"
)

type LockInterface interface {
	LockKey() string
	Execute() error
}

type RedisLockClient struct {
	redisClient *client.Client
	config      *Config
	mu          sync.Mutex
}

const (
	defaultTimeout  = 10 * time.Second
	defaultAttempts = 30
	defaultDelay    = 200 * time.Millisecond
)

func NewRedisLockClient(redisClient *client.Client, opts ...Option) *RedisLockClient {
	//default
	config := &Config{
		timeout:  defaultTimeout,
		attempts: defaultAttempts,
		delay:    defaultDelay,
	}
	for _, opt := range opts {
		opt(config)
	}
	return &RedisLockClient{
		redisClient: redisClient,
		config:      config,
	}
}

func (rlc RedisLockClient) Do(lockInterface LockInterface) error {
	rlc.mu.Lock()
	defer rlc.mu.Unlock()

	var n uint
	for n < rlc.config.attempts {
		lockSrc, err := rlc.redisClient.Get(lockInterface.LockKey())
		if err != nil {
			return err
		}
		if lockSrc == "" || lockSrc == Unlock {
			err = rlc.lock(lockInterface.LockKey())
			if err != nil {
				return err
			}
			err = lockInterface.Execute()
			if err != nil {
				return err
			}
			err := rlc.unlock(lockInterface.LockKey())
			if err != nil {
				return err
			}
		}
		if n == rlc.config.attempts-1 {
			return fmt.Errorf("The number of retries is over")
		}
		time.Sleep(rlc.config.delayType(n, rlc.config))
		n++
		continue
	}
	return fmt.Errorf("unknow error")
}

func (rlc RedisLockClient) lock(key string) error {
	rlc.mu.Lock()
	defer rlc.mu.Unlock()
	return rlc.redisClient.SetByExpiration(key, Lock, rlc.config.timeout)
}

func (rlc RedisLockClient) unlock(key string) error {
	rlc.mu.Lock()
	defer rlc.mu.Unlock()
	return rlc.redisClient.SetByExpiration(key, Unlock, rlc.config.timeout)
}
