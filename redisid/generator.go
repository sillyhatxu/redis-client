package redisid

import (
	"fmt"
	"github.com/sillyhatxu/id-generator"
	"github.com/sillyhatxu/redis-client/redis"
	"time"
)

const (
	defaultGroupLength    int    = 2
	defaultSequenceFormat string = "%04d"
	defaultLifeCycle             = 2 * time.Second
)

type GeneratorClient struct {
	generatorClient *idgenerator.GeneratorClient
}

func NewGeneratorConfig(redisKey string, redisClient *client.Client, opts ...Option) (*GeneratorClient, error) {
	if redisKey == "" {
		return nil, fmt.Errorf("redis key cannot empty")
	}
	if redisClient == nil {
		return nil, fmt.Errorf("redis client is nill")
	}
	i, err := redisClient.Incr(redisKey)
	if err != nil {
		return nil, err
	}
	//default
	config := &Config{
		Prefix:         "",
		GroupLength:    defaultGroupLength,
		SequenceFormat: defaultSequenceFormat,
		LifeCycle:      defaultLifeCycle,
	}
	for _, opt := range opts {
		opt(config)
	}

	return &GeneratorClient{
		generatorClient: idgenerator.NewGeneratorClient(
			redisKey,
			idgenerator.Prefix(config.Prefix),
			idgenerator.GroupLength(config.GroupLength),
			idgenerator.SequenceFormat(config.SequenceFormat),
			idgenerator.Instance(fmt.Sprintf("%d", i)),
			idgenerator.LifeCycle(config.LifeCycle),
		),
	}, nil
}

func (gc GeneratorClient) GeneratorId() (string, error) {
	return gc.GeneratorGroupId("")
}

func (gc GeneratorClient) GeneratorGroupId(src string) (string, error) {
	return gc.generatorClient.GeneratorGroupId(src)
}
