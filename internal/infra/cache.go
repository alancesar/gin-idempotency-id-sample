package infra

import (
	"github.com/allegro/bigcache/v2"
	"github.com/eko/gocache/cache"
	"github.com/eko/gocache/store"
	"time"
)

type BigCacheProvider struct {
	client *cache.Cache
}

func NewBigCacheProvider(config bigcache.Config) *BigCacheProvider {
	bigCacheClient, _ := bigcache.NewBigCache(config)
	bigCacheStore := store.NewBigcache(bigCacheClient, nil)
	return &BigCacheProvider{
		client: cache.New(bigCacheStore),
	}
}

func (b BigCacheProvider) Get(key interface{}) (interface{}, error) {
	return b.client.Get(key)
}

func (b BigCacheProvider) Set(key, value interface{}, ttl time.Duration) error {
	return b.client.Set(key, value, &store.Options{
		Expiration: ttl,
	})
}

func (b BigCacheProvider) Delete(key interface{}) error {
	return b.client.Delete(key)
}
