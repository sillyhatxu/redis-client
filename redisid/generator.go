package redisid

import (
	"fmt"
	"github.com/sillyhatxu/redis-client"
	"sync"
	"time"
)

type LifeCycleType int

const (
	Second LifeCycleType = iota
	Minute
	Hour
)

const (
	defaultGroupLength        int    = 2
	defaultSequenceFormat     string = "%03d"
	defaultLifeCycle                 = Minute
	defaultShards                    = 2048
	defaultLifeWindow                = 24 * time.Hour
	defaultCleanWindow               = 48 * time.Hour
	defaultMaxEntriesInWindow        = 1000 * 10 * 60
	defaultMaxEntrySize              = 500
	defaultVerbose                   = true
	defaultHardMaxCacheSize          = 8192
	defaultHasher                    = ""
)

type GeneratorClient struct {
	redisKey    string
	redisClient *client.Client
	config      *Config
	mu          sync.Mutex
}

func NewGeneratorClient(redisKey string, redisClient *client.Client, opts ...Option) (*GeneratorClient, error) {
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
		redisKey:    redisKey,
		config:      config,
		redisClient: redisClient,
	}, nil
}
func getLifeWindowAndCleanWindow(lifeCycle LifeCycleType) (time.Duration, time.Duration) {
	if lifeCycle == Minute {
		return 1 * time.Minute, 2 * time.Minute
	} else if lifeCycle == Hour {
		return 1 * time.Hour, 2 * time.Hour
	} else {
		return 1 * time.Second, 2 * time.Second
	}
}

func (gc *GeneratorClient) validate() error {
	if gc.redisKey == "" {
		return fmt.Errorf("redis key cannot empty")
	}
	if gc.config == nil {
		return fmt.Errorf("config is nil")
	}
	if gc.redisClient == nil {
		return fmt.Errorf("redis client is nil")
	}
	return nil
}

func (gc *GeneratorClient) GeneratorId() (string, error) {
	return gc.GeneratorGroupId("")
}

func (gc *GeneratorClient) GeneratorGroupId(src string) (string, error) {
	err := gc.validate()
	if err != nil {
		return "", err
	}
	gc.mu.Lock()
	defer gc.mu.Unlock()
	currentTime := time.Now()
	group, err := gc.formatGroup(src)
	if err != nil {
		return "", err
	}
	sequence, err := gc.getSequence(group, currentTime)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s%s", gc.config.Prefix, sequence, gc.getTimeInMillis(currentTime), group), nil
}

func (gc *GeneratorClient) formatGroup(src string) (string, error) {
	if src == "" {
		return "", nil
	}
	hexEncodeSrc := hexEncodeToString(src)
	if len(hexEncodeSrc) > gc.config.GroupLength {
		return hexEncodeSrc[len(hexEncodeSrc)-gc.config.GroupLength:], nil
	} else {
		return hexEncodeSrc, nil
	}
}

func (gc *GeneratorClient) getSequence(group string, currentTime time.Time) (string, error) {
	key := fmt.Sprintf("%s_%s_%s", gc.redisKey, group, gc.getKeySuffix(currentTime))
	sequence, err := gc.redisClient.IncrByExpiration(key, gc.getExpiration())
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(gc.config.SequenceFormat, sequence), nil
}

func (gc *GeneratorClient) getExpiration() time.Duration {
	if gc.config.LifeCycle == Minute {
		return time.Minute
	} else if gc.config.LifeCycle == Hour {
		return time.Hour
	} else {
		return time.Second
	}
}

func (gc *GeneratorClient) getKeySuffix(currentTime time.Time) string {
	hr, min, sec := currentTime.Clock()
	if gc.config.LifeCycle == Minute {
		return fmt.Sprintf("%d_%d", hr, min)
	} else if gc.config.LifeCycle == Hour {
		return fmt.Sprintf("%d", hr)
	} else {
		return fmt.Sprintf("%d_%d_%d", hr, min, sec)
	}
}

func (gc *GeneratorClient) getTimeInMillis(currentTime time.Time) string {
	return Int2String(currentTime.Unix() / getLifeCycleNumber(gc.config.LifeCycle))
}
