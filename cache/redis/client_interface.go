package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisClient is the underlying redis client interface, go-redis library
// we need this interface mainly to unify implementation of cluster and single instance mode
type redisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	MGet(ctx context.Context, keys ...string) *redis.SliceCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Pipeline() redis.Pipeliner
	Watch(ctx context.Context, fn func(tx *redis.Tx) error, keys ...string) error

	Ping(ctx context.Context) *redis.StatusCmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
	Close() error

	redisClientForTest
}

// redisClientForTest lists functions used only in unit tests
type redisClientForTest interface {
	TTL(ctx context.Context, key string) *redis.DurationCmd
}
