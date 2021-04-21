package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	i, err := redisInit(t)
	assert.Nil(t, err)
	assert.NotNil(t, i)
}

func TestRedis_Set(t *testing.T) {
	r, err := redisInit(t)
	assert.Nil(t, err)

	testData := struct {
		Field string
		F     float64
	}{"1", 200.1}

	err = r.Set("test", testData, time.Second)
	assert.Nil(t, err)

	var newValue struct {
		Field string
		F     float64
	}
	err = r.Get("test", &newValue)
	assert.Nil(t, err)
	assert.Equal(t, testData, newValue)

	ttl := r.client.TTL("test")
	assert.Equal(t, time.Second, ttl.Val())
}

func TestRedis_Get(t *testing.T) {
	r, err := redisInit(t)
	assert.Nil(t, err)

	type TestStruct struct {
		Field string
		F     float64
	}
	testData := TestStruct{"1", 200.1}

	err = r.Set("test", testData, time.Second)
	assert.Nil(t, err)

	var newValue TestStruct
	err = r.Get("test", &newValue)
	assert.Nil(t, err)
	assert.Equal(t, testData, newValue)

	ttl := r.client.TTL("test")
	assert.Equal(t, time.Second, ttl.Val())

	var empty interface{}
	err = r.Get("1", empty)
	assert.NotNil(t, err)
	assert.Nil(t, empty)
}

func TestRedis_Delete(t *testing.T) {
	r, err := redisInit(t)
	assert.Nil(t, err)
	err = r.Set("test", []byte{0, 1}, time.Second)
	assert.Nil(t, err)

	err = r.Delete("test")
	assert.Nil(t, err)

	var v []byte
	err = r.Get("test", &v)
	assert.NotNil(t, err)
	assert.Equal(t, string([]byte{}), string(v))
}

func TestRedis_IsAvailable(t *testing.T) {
	r, err := redisInit(t)
	assert.Nil(t, err)

	assert.True(t, r.IsAvailable())
}

func TestRedis_Reconnect(t *testing.T) {
	r, err := redisInit(t)
	assert.Nil(t, err)

	mr, err := miniredis.Run()
	assert.NotNil(t, mr)
	assert.Nil(t, err)

	err = r.Reconnect(fmt.Sprintf("redis://%s", mr.Addr()))
	assert.NoError(t, err)
}

func redisInit(t *testing.T) (*Redis, error) {
	mr, err := miniredis.Run()
	assert.NotNil(t, mr)
	assert.Nil(t, err)

	c, err := Init(fmt.Sprintf("redis://%s", mr.Addr()), time.Minute)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	return c, nil
}
