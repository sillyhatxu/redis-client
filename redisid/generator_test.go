package redisid

import (
	"fmt"
	"github.com/sillyhatxu/redis-client/redis"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

var redisClient *client.Client

func init() {
	rClient, err := client.NewRedisClient("127.0.0.1:6379", "", 0)
	if err != nil {
		panic(err)
	}
	redisClient = rClient
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now().UnixNano())
	fmt.Println(time.Now().UnixNano() / int64(time.Millisecond))
	fmt.Println(strconv.FormatInt(time.Now().Unix(), 10))
}

func TestGeneratorId(t *testing.T) {
	generatorClient, err := NewGeneratorClient("id.generator.seq.test", redisClient)
	assert.Nil(t, err)
	id, err := generatorClient.GeneratorId()
	assert.Nil(t, err)
	fmt.Println(id)
}

func TestGeneratorGroupId(t *testing.T) {
	generatorClient, err := NewGeneratorClient("id.generator.seq.test.group", redisClient, Prefix("T"))
	assert.Nil(t, err)
	for i := 0; i < 50; i++ {
		id, err := generatorClient.GeneratorGroupId("test")
		assert.Nil(t, err)
		fmt.Println(id)
	}
}
