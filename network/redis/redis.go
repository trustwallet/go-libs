package redis

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

var NotFound = errors.New("not found")

type Cache struct {
	client redis.Client
}

func Init(host string) (Cache, error) {
	options, err := redis.ParseURL(host)
	if err != nil {
		return Cache{}, err
	}
	client := redis.NewClient(options)
	return Cache{client: *client}, nil
}

func (c Cache) Get(key string) ([]byte, error) {
	cmd := c.client.Get(key)
	if cmd.Err() == redis.Nil {
		return nil, NotFound
	} else if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return []byte(cmd.Val()), nil
}

func (c Cache) Set(key string, value []byte, expiration time.Duration) error {
	cmd := c.client.Set(key, value, expiration)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (c Cache) Delete(key string) error {
	cmd := c.client.Del(key)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
