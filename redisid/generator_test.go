package redisid

import (
	"fmt"
	"github.com/sillyhatxu/go-utils/hashset"
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
	var generatorConfig = NewGeneratorConfig("id.generator.seq.test", redisClient)
	id, err := generatorConfig.GeneratorId()
	assert.Nil(t, err)
	fmt.Println(id)
}

var testHashSet = hashset.New()

func TestGeneratorGroupId(t *testing.T) {
	var generatorConfig = NewGeneratorConfig("id.generator.seq.test.group", redisClient, Prefix("T"))
	go func() {
		for i := 0; i < 50; i++ {
			id, err := generatorConfig.GeneratorGroupId("test")
			assert.Nil(t, err)
			fmt.Println(id)
			testHashSet.Add(id)
			fmt.Println(testHashSet.Size())
		}
	}()
	go func() {
		for i := 0; i < 50; i++ {
			id, err := generatorConfig.GeneratorGroupId("test")
			assert.Nil(t, err)
			fmt.Println(id)
			testHashSet.Add(id)
			fmt.Println(testHashSet.Size())
		}
	}()
	go func() {
		for i := 0; i < 50; i++ {
			id, err := generatorConfig.GeneratorGroupId("test")
			assert.Nil(t, err)
			fmt.Println(id)
			testHashSet.Add(id)
		}
		fmt.Println(testHashSet.Size())
	}()
	time.Sleep(30 * time.Second)
}
