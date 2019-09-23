package redisid

import (
	"github.com/sillyhatxu/id-generator"
)

type Config struct {
	Prefix         string
	GroupLength    int
	SequenceFormat string
	LifeCycle      idgenerator.LifeCycleType
}

type Option func(*Config)

func Prefix(prefix string) Option {
	return func(c *Config) {
		c.Prefix = prefix
	}
}

func GroupLength(groupLength int) Option {
	return func(c *Config) {
		c.GroupLength = groupLength
	}
}

func SequenceFormat(sequenceFormat string) Option {
	return func(c *Config) {
		c.SequenceFormat = sequenceFormat
	}
}

func LifeCycle(lifeCycle idgenerator.LifeCycleType) Option {
	return func(c *Config) {
		c.LifeCycle = lifeCycle
	}
}
