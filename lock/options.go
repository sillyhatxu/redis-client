package redislock

import (
	"time"
)

type DelayTypeFunc func(n uint, config *Config) time.Duration

type Config struct {
	timeout   time.Duration
	attempts  uint
	delay     time.Duration
	delayType DelayTypeFunc
}

type Option func(*Config)

func Timeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.timeout = timeout
	}
}

func Attempts(attempts uint) Option {
	return func(c *Config) {
		c.attempts = attempts
	}
}

func Delay(delay time.Duration) Option {
	return func(c *Config) {
		c.delay = delay
	}
}

func DelayType(delayType DelayTypeFunc) Option {
	return func(c *Config) {
		c.delayType = delayType
	}
}

func BackOffDelay(n uint, config *Config) time.Duration {
	return config.delay * (1 << n)
}
