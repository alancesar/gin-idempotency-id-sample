package cache

import (
	"encoding/json"
	"net/http"
	"time"
)

type (
	IdempotencyCache struct {
		provider Provider
		ttl      time.Duration
	}

	Data struct {
		Headers     http.Header
		ContentType string
		Body        []byte
		StatusCode  int
	}

	key struct {
		IdempotencyID, URL string
	}

	Provider interface {
		Get(key interface{}) (interface{}, error)
		Set(key, value interface{}, ttl time.Duration) error
	}
)

func NewIdempotencyCache(cache Provider, ttl time.Duration) *IdempotencyCache {
	return &IdempotencyCache{
		provider: cache,
		ttl:      ttl,
	}
}

func (c IdempotencyCache) Set(idempotencyID, url string, data Data) error {
	body, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	key := key{
		IdempotencyID: idempotencyID,
		URL:           url,
	}
	return c.provider.Set(key, body, c.ttl)
}

func (c IdempotencyCache) Get(idempotencyID, url string) (Data, error) {
	key := key{
		IdempotencyID: idempotencyID,
		URL:           url,
	}
	data, err := c.provider.Get(key)
	if err != nil {
		return Data{}, err
	}

	var output Data
	err = json.Unmarshal(data.([]uint8), &output)
	return output, err
}
