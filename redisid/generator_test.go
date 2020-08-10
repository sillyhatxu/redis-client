package redisid

import (
	client "github.com/sillyhatxu/redis-client"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

var redisClient *client.Client
var once sync.Once

const count = 1000

func setup() {
	//rc, err := client.NewRedisClient("127.0.0.1:16379", "sillyhatpassword", 0)
	rc, err := client.NewRedisClient("127.0.0.1:6379", "", 0)
	if err != nil {
		panic(err)
	}
	redisClient = rc
}

func TestGeneratorClient_GeneratorGroupIdDefault(t *testing.T) {
	once.Do(setup)
	gc, err := NewGeneratorClient("test", redisClient, Prefix("SH"))
	assert.Nil(t, err)
	check := make(map[string]bool)
	for i := 0; i < count; i++ {
		id, err := gc.GeneratorId()
		assert.Nil(t, err)
		check[id] = true
	}
	for i := 0; i < count; i++ {
		id, err := gc.GeneratorGroupId("TEST_ID")
		assert.Nil(t, err)
		check[id] = true
	}
	assert.EqualValues(t, count*2, len(check))
}

func TestGeneratorClient_GeneratorGroupId(t *testing.T) {
	once.Do(setup)
	gc, err := NewGeneratorClient("test", redisClient, Prefix("SH"), GroupLength(3), SequenceFormat("%03d"))
	assert.Nil(t, err)
	check := make(map[string]bool)
	for i := 0; i < count; i++ {
		id, err := gc.GeneratorGroupId("TEST_ID")
		assert.Nil(t, err)
		check[id] = true
	}
	assert.EqualValues(t, count, len(check))
}

func TestGeneratorClient_GeneratorGroupIdLifeCycle(t *testing.T) {
	once.Do(setup)
	gc, err := NewGeneratorClient("test", redisClient, Prefix("SH"), GroupLength(3), SequenceFormat("%03d"), LifeCycle(Second))
	assert.Nil(t, err)
	check := make(map[string]bool)
	//var test []string
	for i := 0; i < count; i++ {
		id, err := gc.GeneratorGroupId("TEST_ID")
		assert.Nil(t, err)
		check[id] = true
		//test = append(test, id)
	}
	//for j, t2 := range test {
	//	fmt.Println(j, t2)
	//}
	//if count != len(check) {
	//	for i, t1 := range test {
	//		for j, t2 := range test {
	//			if i != j && t1 == t2 {
	//				fmt.Println(t1)
	//				break
	//			}
	//		}
	//	}
	//}
	assert.EqualValues(t, count, len(check))
}
