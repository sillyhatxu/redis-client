package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	redisClient, err := NewRedisClient("127.0.0.1:16379", "sillyhatpassword", 0)
	assert.Nil(t, err)
	test, err := redisClient.Get("test-src-1")
	assert.Nil(t, err)
	assert.EqualValues(t, test, "")
}

func TestSetGet(t *testing.T) {
	redisClient, err := NewRedisClient("127.0.0.1:16379", "sillyhatpassword", 0)
	assert.Nil(t, err)
	err = redisClient.SetByExpiration("test-src", "1", time.Second)
	assert.Nil(t, err)
	test, err := redisClient.Get("test-src")
	assert.Nil(t, err)
	assert.EqualValues(t, "1", test)
	time.Sleep(time.Second)
	test, err = redisClient.Get("test-src")
	assert.Nil(t, err)
	assert.EqualValues(t, "", test)
}
