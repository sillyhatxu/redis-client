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
	for i := 0; i < 100; i++ {
		id, err := gc.GeneratorGroupId("TEST_ID")
		assert.Nil(t, err)
		fmt.Println(id)
	}
}

func TestGeneratorClient_GeneratorGroupId(t *testing.T) {
	redisClient, err := client.NewRedisClient("127.0.0.1:16379", "sillyhatpassword", 0)
	gc, err := NewGeneratorClient("test", redisClient, Prefix("SH"), GroupLength(3), SequenceFormat("%03d"))
	assert.Nil(t, err)
	for i := 0; i < 100; i++ {
		id, err := gc.GeneratorGroupId("group")
		assert.Nil(t, err)
		fmt.Println(id)
	}
}

func TestGeneratorClient_GeneratorGroupIdInstance(t *testing.T) {
	redisClient, err := client.NewRedisClient("127.0.0.1:16379", "sillyhatpassword", 0)
	gc, err := NewGeneratorClient(
		"test", redisClient,
		Prefix("GT"),
		GroupLength(3),
		SequenceFormat("%03d"),
		LifeCycle(Minute),
	)
	assert.Nil(t, err)
	for i := 0; i < 500; i++ {
		id, err := gc.GeneratorGroupId("group")
		assert.Nil(t, err)
		fmt.Println(id)
		time.Sleep(500 * time.Millisecond)
		//GT826176437022516
	}
}
