package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrNotFound = fmt.Errorf("not found")

const healthCheckTimeout = 2 * time.Second

type Redis struct {
	client    redisClient
	isCluster bool
}

type Option func(o *redis.Options)
type ClusterOption func(o *redis.ClusterOptions)

// WithTLS option allows to set TLS config to initiate
// connection to Redis host with TLS transport protocol enabled.
// Example configuration for tests, self-signed untrusted certs.
//
//	&tls.Config{
//		InsecureSkipVerify: true,
//	}
//
// Example of configuration for production usage
//
//	&tls.Config{
//		MinVersion: tls.VersionTLS12,
//	}
func WithTLS(cfg *tls.Config) Option {
	return func(o *redis.Options) {
		o.TLSConfig = cfg
	}
}

func WithClusterTLS(cfg *tls.Config) ClusterOption {
	return func(o *redis.ClusterOptions) {
		o.TLSConfig = cfg
	}
}

func InitWithTLS(url string, secure bool, opts ...Option) (*Redis, error) {
	var options []Option
	if secure {
		options = append(options, WithTLS(&tls.Config{InsecureSkipVerify: true}))
	}
	options = append(options, opts...)

	redis, err := Init(context.Background(), url, options...)
	if err != nil {
		return nil, err
	}
	return redis, nil
}

func Init(ctx context.Context, host string, opts ...Option) (*Redis, error) {
	options, err := redis.ParseURL(host)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(options)
	}

	client := redis.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{client: client}, nil
}

func InitClusterWithTLS(url string, secure bool, opts ...ClusterOption) (*Redis, error) {
	var options []ClusterOption
	if secure {
		options = append(options, WithClusterTLS(&tls.Config{InsecureSkipVerify: true}))
	}
	options = append(options, opts...)

	redis, err := InitCluster(context.Background(), url, options...)
	if err != nil {
		return nil, err
	}
	return redis, nil
}

// InitCluster inits redis connection with cluster mode
//
// In cluster mode, redisURL supports multiple addresses. Example:
//
//	redis://user:password@localhost:6789?dial_timeout=3&read_timeout=6s&addr=localhost:6790&addr=localhost:6791
//	is equivalent to:
//	&ClusterOptions{
//		Addr:        ["localhost:6789", "localhost:6790", "localhost:6791"]
//		DialTimeout: 3 * time.Second, // no time unit = seconds
//		ReadTimeout: 6 * time.Second,
//	}
func InitCluster(ctx context.Context, redisURL string, opts ...ClusterOption) (*Redis, error) {
	options, err := redis.ParseClusterURL(redisURL)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(options)
	}

	client := redis.NewClusterClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{client: client, isCluster: true}, nil
}

func (r *Redis) Get(ctx context.Context, key string, receiver interface{}) error {
	cmd := r.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redis.Nil) {
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

func (r *Redis) GetBytes(ctx context.Context, key string) ([]byte, error) {
	cmd := r.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redis.Nil) {
		return nil, ErrNotFound
	} else if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return cmd.Bytes()
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

// Scan return keys by pattern
func (r *Redis) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, error) {
	iter := r.client.Scan(ctx, cursor, match, count).Iterator()

	res := make([]string, 0)
	for iter.Next(ctx) {
		res = append(res, iter.Val())
	}

	return res, iter.Err()
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

func (r *Redis) SetBytes(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	cmd := r.client.Set(ctx, key, value, expiration)
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

func (r *Redis) Watch(ctx context.Context, fn func(tx *redis.Tx) error, keys ...string) error {
	return r.client.Watch(ctx, fn, keys...)
}

func (r *Redis) IsAvailable(ctx context.Context) bool {
	return r.client.Ping(ctx).Err() == nil
}

func (r *Redis) Reconnect(ctx context.Context, host string) error {
	if r.isCluster {
		return r.reconnectCluster(ctx, host)
	}

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

func (r *Redis) reconnectCluster(ctx context.Context, redisURL string) error {
	options, err := redis.ParseClusterURL(redisURL)
	if err != nil {
		return err
	}

	client := redis.NewClusterClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return err
	}

	r.client = client
	if err := r.client.Ping(ctx).Err(); err != nil {
		return err
	}

	return nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) HealthCheck() error {
	if r == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
	defer cancel()

	if !r.IsAvailable(ctx) {
		return errors.New("redis is not available")
	}

	return nil
}
