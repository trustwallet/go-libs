package redis

import (
	"context"
	"encoding/json"
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

	err = r.Set(context.TODO(), "test", testData, time.Second)
	assert.Nil(t, err)

	var newValue struct {
		Field string
		F     float64
	}
	err = r.Get(context.TODO(), "test", &newValue)
	assert.Nil(t, err)
	assert.Equal(t, testData, newValue)

	ttl := r.client.TTL(context.TODO(), "test")
	assert.Equal(t, time.Second, ttl.Val())
}

func TestRedis_MSet_MGet(t *testing.T) {
	r, err := redisInit(t)
	assert.Nil(t, err)

	type testDataType struct {
		Field string
		F     float64
	}
	testData := map[string]interface{}{
		"test1": testDataType{"1", 200.1},
		"test2": testDataType{"1", 200.1},
	}

	err = r.MSet(context.TODO(), testData, time.Second)
	assert.Nil(t, err)

	keys := []string{"test1", "test2"}
	values, err := r.MGet(context.TODO(), keys...)
	assert.Nil(t, err)
	assert.Equal(t, len(keys), len(values))

	for i, key := range keys {
		var data testDataType
		err := json.Unmarshal(values[i], &data)
		assert.NoError(t, err)
		assert.Equal(t, testData[key], data)
	}

	ttl := r.client.TTL(context.TODO(), "test1")
	assert.Equal(t, time.Second, ttl.Val())
	ttl = r.client.TTL(context.TODO(), "test2")
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

	err = r.Set(context.TODO(), "test", testData, time.Second)
	assert.Nil(t, err)

	var newValue TestStruct
	err = r.Get(context.TODO(), "test", &newValue)
	assert.Nil(t, err)
	assert.Equal(t, testData, newValue)

	ttl := r.client.TTL(context.TODO(), "test")
	assert.Equal(t, time.Second, ttl.Val())

	var empty interface{}
	err = r.Get(context.TODO(), "1", empty)
	assert.Equal(t, err, ErrNotFound)
	assert.Nil(t, empty)
}

func TestRedis_Delete(t *testing.T) {
	r, err := redisInit(t)
	assert.Nil(t, err)
	err = r.Set(context.TODO(), "test", []byte{0, 1}, time.Second)
	assert.Nil(t, err)

	err = r.Delete(context.TODO(), "test")
	assert.Nil(t, err)

	var v []byte
	err = r.Get(context.TODO(), "test", &v)
	assert.NotNil(t, err)
	assert.Equal(t, string([]byte{}), string(v))
}

func TestRedis_IsAvailable(t *testing.T) {
	r, err := redisInit(t)
	assert.Nil(t, err)

	assert.True(t, r.IsAvailable(context.TODO()))
}

func TestRedis_Reconnect(t *testing.T) {
	r, err := redisInit(t)
	assert.Nil(t, err)

	mr, err := miniredis.Run()
	assert.NotNil(t, mr)
	assert.Nil(t, err)

	err = r.Reconnect(context.TODO(), fmt.Sprintf("redis://%s", mr.Addr()))
	assert.NoError(t, err)
}

func redisInit(t *testing.T) (*Redis, error) {
	mr, err := miniredis.Run()
	assert.NotNil(t, mr)
	assert.Nil(t, err)

	c, err := Init(context.TODO(), fmt.Sprintf("redis://%s", mr.Addr()))
	assert.Nil(t, err)
	assert.NotNil(t, c)

	return c, nil
}
