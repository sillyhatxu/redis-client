package redisid

import (
	"fmt"
	client "github.com/sillyhatxu/redis-client"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGeneratorClient_GeneratorGroupIdDefault(t *testing.T) {
	redisClient, err := client.NewRedisClient("127.0.0.1:16379", "sillyhatpassword", 0)
	gc, err := NewGeneratorClient("test", redisClient, Prefix("SH"))
	assert.Nil(t, err)
	for i := 0; i < 10; i++ {
		id, err := gc.GeneratorId()
		assert.Nil(t, err)
		fmt.Println(id)
	}
	for i := 0; i < 10; i++ {
		id, err := gc.GeneratorGroupId("TEST_ID")
		assert.Nil(t, err)
		fmt.Println(id)
	}
}

func TestGeneratorClient_GeneratorGroupId(t *testing.T) {
	redisClient, err := client.NewRedisClient("127.0.0.1:16379", "sillyhatpassword", 0)
	gc, err := NewGeneratorClient("test", redisClient, Prefix("SH"), GroupLength(3), SequenceFormat("%03d"))
	assert.Nil(t, err)
	for i := 0; i < 10; i++ {
		id, err := gc.GeneratorGroupId("TEST_ID")
		assert.Nil(t, err)
		fmt.Println(id)
	}
}

func TestGeneratorClient_GeneratorGroupIdLifeCycle(t *testing.T) {
	redisClient, err := client.NewRedisClient("127.0.0.1:16379", "sillyhatpassword", 0)
	gc, err := NewGeneratorClient("test", redisClient, Prefix("SH"), GroupLength(3), SequenceFormat("%03d"), LifeCycle(Second))
	assert.Nil(t, err)
	for i := 0; i < 500; i++ {
		id, err := gc.GeneratorGroupId("TEST_ID")
		assert.Nil(t, err)
		fmt.Println(id)
		time.Sleep(50 * time.Millisecond)
	}
}
