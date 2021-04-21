package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

var ErrNotFound = fmt.Errorf("not found")

func Init(host string) (*Redis, error) {
	options, err := redis.ParseURL(host)

	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &Redis{client: client}, nil
}

func (r *Redis) Get(key string, receiver interface{}) error {
	cmd := r.client.Get(key)
	if cmd.Err() == redis.Nil {
		return ErrNotFound
	} else if cmd.Err() != nil {
		return cmd.Err()
	}

	err := json.Unmarshal([]byte(cmd.Val()), receiver)
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	cmd := r.client.Set(key, data, expiration)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *Redis) Delete(key string) error {
	cmd := r.client.Del(key)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *Redis) IsAvailable() bool {
	return r.client.Ping().Err() == nil
}

func (r *Redis) Reconnect(host string) error {
	options, err := redis.ParseURL(host)
	if err != nil {
		return err
	}

	client := redis.NewClient(options)
	if err := client.Ping().Err(); err != nil {
		return err
	}

	r.client = client
	if err := r.client.Ping().Err(); err != nil {
		return err
	}

	return nil
}
