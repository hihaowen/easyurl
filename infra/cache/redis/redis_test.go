package redis

import (
	"github.com/bmizerany/assert"
	"log"
	"testing"
)

func TestRedis(t *testing.T) {
	// connect
	Connect(RedisConf{"127.0.0.1:6379", ""})

	key := "testKey1"
	value := "testVal1"

	// set
	Set(key, value)

	// get
	getValue, err := Get(key)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, value, string(getValue))
}
