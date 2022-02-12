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

	Provider interface {
		Get(key interface{}) (interface{}, error)
		Set(key, value interface{}, ttl time.Duration) error
	}
)

func (d Data) WriteHeaders(writer http.ResponseWriter) {
	for k, v := range d.Headers {
		writer.Header().Set(k, v[0])
	}
}

func NewIdempotencyCache(cache Provider, ttl time.Duration) *IdempotencyCache {
	return &IdempotencyCache{
		provider: cache,
		ttl:      ttl,
	}
}

func (c IdempotencyCache) Set(key interface{}, data Data) error {
	body, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	return c.provider.Set(key, body, c.ttl)
}

func (c IdempotencyCache) Get(key interface{}) (Data, error) {
	data, err := c.provider.Get(key)
	if err != nil {
		return Data{}, err
	}

	var output Data
	err = json.Unmarshal(data.([]uint8), &output)
	return output, err
}
