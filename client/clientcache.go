package client

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

var memoryCache *memCache

func init() {
	memoryCache = &memCache{cache: cache.New(5*time.Minute, 5*time.Minute)}
}

type memCache struct {
	cache *cache.Cache
}

func (r *Request) PostWithCache(result interface{}, path string, body interface{}, cache time.Duration) error {
	key := r.generateKey(path, nil, body)
	err := memoryCache.getCache(key, result)
	if err == nil {
		return nil
	}

	err = r.Post(result, path, body)
	if err != nil {
		return err
	}

	return memoryCache.setCache(key, result, cache)
}

func (r *Request) PostWithCacheAndContext(result interface{}, path string, body interface{}, cache time.Duration, ctx context.Context) error {
	key := r.generateKey(path, nil, body)
	err := memoryCache.getCache(key, result)
	if err == nil {
		return nil
	}

	err = r.PostWithContext(result, path, body, ctx)
	if err != nil {
		return err
	}

	return memoryCache.setCache(key, result, cache)
}

func (r *Request) GetWithCache(result interface{}, path string, query url.Values, cache time.Duration) error {
	key := r.generateKey(path, query, nil)
	err := memoryCache.getCache(key, result)
	if err == nil {
		return nil
	}

	err = r.Get(result, path, query)
	if err != nil {
		return err
	}

	return memoryCache.setCache(key, result, cache)
}

func (r *Request) GetWithCacheAndContext(result interface{}, path string, query url.Values, cache time.Duration, ctx context.Context) error {
	key := r.generateKey(path, query, nil)
	err := memoryCache.getCache(key, result)
	if err == nil {
		return nil
	}

	err = r.Get(result, path, query)
	if err != nil {
		return err
	}

	return memoryCache.setCache(key, result, cache)
}

func (mc *memCache) setCache(key string, value interface{}, duration time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return errors.New(err.Error() + " client cache cannot marshal cache object")
	}
	memoryCache.cache.Set(key, b, duration)
	return nil
}

func (mc *memCache) getCache(key string, value interface{}) error {
	c, ok := mc.cache.Get(key)
	if !ok {
		return errors.New("validator cache: invalid cache key")
	}
	r, ok := c.([]byte)
	if !ok {
		return errors.New("validator cache: failed to cast cache to bytes")
	}
	err := json.Unmarshal(r, value)
	if err != nil {
		return errors.New(err.Error() + " not found")
	}
	return nil
}

func (r *Request) generateKey(path string, query url.Values, body interface{}) string {
	var queryStr = ""
	if query != nil {
		queryStr = query.Encode()
	}
	requestUrl := strings.Join([]string{r.GetBase(path), queryStr}, "?")
	var b []byte
	if body != nil {
		b, _ = json.Marshal(body)
	}
	hash := sha1.Sum(append([]byte(requestUrl), b...))
	return base64.URLEncoding.EncodeToString(hash[:])
}
