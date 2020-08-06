package client

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

type Client struct {
	address     string
	password    string
	db          int
	redisClient *redis.Client
}

func NewRedisClient(address, password string, db int) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	if client == nil {
		return nil, fmt.Errorf("new redis client error. address : %s; password : %s;db : %d", address, password, db)
	}
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &Client{
		address:     address,
		password:    password,
		db:          db,
		redisClient: client,
	}, nil
}

func (c *Client) GetClient() (*redis.Client, error) {
	_, err := c.redisClient.Ping().Result()
	if err != nil {
		client := redis.NewClient(&redis.Options{
			Addr:     c.address,
			Password: c.password,
			DB:       c.db,
		})
		if client == nil {
			return nil, fmt.Errorf("new redis client error. address : %s; password : %s;db : %d", c.address, c.password, c.db)
		}
		_, err := client.Ping().Result()
		if err != nil {
			return nil, err
		}
		c.redisClient = client
	}
	return c.redisClient, nil
}

func (c *Client) Get(key string) (string, error) {
	client, err := c.GetClient()
	if err != nil {
		return "", err
	}
	value, err := client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (c *Client) Set(key, value string) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}
	err = client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SetByExpiration(key, value string, expiration time.Duration) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}
	err = client.Set(key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Exists(key string) (bool, error) {
	client, err := c.GetClient()
	if err != nil {
		return false, err
	}
	count := client.Exists(key).Val()
	return count > 0, nil
}

func (c *Client) Expire(key string, expiration time.Duration) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}
	err = client.Expire(key, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Incr(key string) (int64, error) {
	client, err := c.GetClient()
	if err != nil {
		return 0, err
	}
	seq := client.Incr(key).Val()
	return seq, nil
}

func (c *Client) IncrByExpiration(key string, expiration time.Duration) (result int64, err error) {
	result, err = c.Incr(key)
	if err != nil {
		return
	}
	err = c.Expire(key, expiration)
	if err != nil {
		return
	}
	return
}
