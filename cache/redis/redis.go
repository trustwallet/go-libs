package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

var ErrNotFound = fmt.Errorf("not found")

func Init(ctx context.Context, host string) (*Redis, error) {
	options, err := redis.ParseURL(host)

	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{client: client}, nil
}

func (r *Redis) Get(ctx context.Context, key string, receiver interface{}) error {
	cmd := r.client.Get(ctx, key)
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

// MGet returns slice with length == len(key)
// Resulting slice's item is nil if there is no value in cache
func (r *Redis) MGet(ctx context.Context, key ...string) ([][]byte, error) {
	cmd := r.client.MGet(ctx, key...)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	result := make([][]byte, len(cmd.Val()))
	for i, v := range cmd.Val() {
		if v != nil {
			result[i] = []byte(v.(string))
		}
	}

	return result, nil
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	cmd := r.client.Set(ctx, key, data, expiration)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *Redis) MSet(ctx context.Context, pairs map[string]interface{}, expiration time.Duration) error {
	p := r.client.Pipeline()

	for k, v := range pairs {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		cmd := p.Set(ctx, k, data, expiration)
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}

	_, err := p.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	cmd := r.client.Del(ctx, key)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *Redis) IsAvailable(ctx context.Context) bool {
	return r.client.Ping(ctx).Err() == nil
}

func (r *Redis) Reconnect(ctx context.Context, host string) error {
	options, err := redis.ParseURL(host)
	if err != nil {
		return err
	}

	client := redis.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return err
	}

	r.client = client
	if err := r.client.Ping(ctx).Err(); err != nil {
		return err
	}

	return nil
}
