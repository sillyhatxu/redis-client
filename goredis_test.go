package client

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

var redisClient *Client
var once sync.Once

func setup() {
	//rc, err := NewRedisClient("127.0.0.1:16379", "sillyhatpassword", 0)
	rc, err := NewRedisClient("127.0.0.1:6379", "", 0)
	if err != nil {
		panic(err)
	}
	redisClient = rc
}

func TestGet(t *testing.T) {
	once.Do(setup)
	test, err := redisClient.Get("test-src-1")
	assert.Nil(t, err)
	assert.EqualValues(t, test, "")
}

func TestSetGet(t *testing.T) {
	once.Do(setup)
	err := redisClient.SetByExpiration("test-src", "1", time.Second)
	assert.Nil(t, err)
	test, err := redisClient.Get("test-src")
	assert.Nil(t, err)
	assert.EqualValues(t, "1", test)
	time.Sleep(time.Second)
	test, err = redisClient.Get("test-src")
	assert.Nil(t, err)
	assert.EqualValues(t, "", test)
}
