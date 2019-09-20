package redisid

import (
	"fmt"
	"github.com/sillyhatxu/redis-client/redis"
	"hash/fnv"
	"strconv"
	"sync"
	"time"
)

const (
	defaultAddress        string = "127.0.0.1:6379"
	defaultPassword       string = ""
	defaultDB             int    = 0
	defaultGroupLength    int    = 2
	defaultSequenceFormat string = "%04d"
	defaultLifeCycle             = 5 * time.Second
)

type GeneratorConfig struct {
	redisKey string
	config   *Config
	mu       sync.Mutex
}

func (gc GeneratorConfig) validate() error {
	if gc.redisKey == "" {
		return fmt.Errorf("redis key cannot empty")
	}
	if gc.config == nil {
		return fmt.Errorf("redis config is nill")
	}
	if gc.config.RedisClient == nil {
		return fmt.Errorf("redis client is nill")
	}
	return nil
}

func NewGeneratorConfig(redisKey string, redisClient *client.Client, opts ...Option) *GeneratorConfig {
	//default
	config := &Config{
		RedisClient:    redisClient,
		Prefix:         "",
		GroupLength:    defaultGroupLength,
		SequenceFormat: defaultSequenceFormat,
		LifeCycle:      defaultLifeCycle,
	}
	for _, opt := range opts {
		opt(config)
	}
	return &GeneratorConfig{
		redisKey: redisKey,
		config:   config,
	}
}

func (gc GeneratorConfig) GeneratorId() (string, error) {
	err := gc.validate()
	if err != nil {
		return "", err
	}
	gc.mu.Lock()
	defer gc.mu.Unlock()
	sequence, err := gc.getRedisSequence("")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s", gc.config.Prefix, getTimeInMillis(), sequence), nil
}

func (gc GeneratorConfig) GeneratorGroupId(userId string) (string, error) {
	err := gc.validate()
	if err != nil {
		return "", err
	}
	gc.mu.Lock()
	defer gc.mu.Unlock()
	group, err := gc.formatGroup(userId)
	if err != nil {
		return "", err
	}
	sequence, err := gc.getRedisSequence(group)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s%s", gc.config.Prefix, getTimeInMillis(), sequence, group), nil
}

func (gc GeneratorConfig) formatGroup(src string) (string, error) {
	hashSrc, err := hash(src)
	if err != nil {
		return "", err
	}
	formatUintSrc := strconv.FormatUint(hashSrc, 10)
	return string(formatUintSrc[len(formatUintSrc)-gc.config.GroupLength:]), nil
}

func (gc GeneratorConfig) getRedisSequence(group string) (string, error) {
	redisKey := fmt.Sprintf("%s%s", gc.redisKey, group)
	seq, err := gc.config.RedisClient.Incr(redisKey)
	if seq == 1 {
		err := gc.config.RedisClient.Expire(redisKey, gc.config.LifeCycle)
		if err != nil {
			return "", err
		}
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(gc.config.SequenceFormat, seq), nil
}

func hash(s string) (uint64, error) {
	h := fnv.New64a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}

func getTimeInMillis() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
