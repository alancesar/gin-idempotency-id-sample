package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"idempotency/internal/api/writter"
	"idempotency/internal/cache"
	"net/http"
)

const (
	idempotencyIDKey = "Idempotency-Id"
)

type (
	Cache interface {
		Get(key interface{}) (cache.Data, error)
		Set(key interface{}, data cache.Data) error
		Lock(key interface{}) error
		Unlock(key interface{}) error
	}

	Opts struct {
		CacheKeyGenFn   CacheKeyGenFn
		CacheCriteriaFn CacheCriteriaFn
	}

	CacheKeyGenFn   func(r *http.Request) interface{}
	CacheCriteriaFn func(writer gin.ResponseWriter) bool
)

func KeyAndURL(r *http.Request) interface{} {
	return struct {
		Key string
		URL string
	}{
		Key: r.Header.Get(idempotencyIDKey),
		URL: r.RequestURI,
	}
}

func CacheOnlySuccess(w gin.ResponseWriter) bool {
	return w.Status() >= http.StatusOK && w.Status() < http.StatusBadRequest
}

func Idempotency(cache Cache, opts ...Opts) gin.HandlerFunc {
	keyFn, cacheCriteriaFn := setOpts(opts...)

	return func(ctx *gin.Context) {
		if !isSomeHTTPMethod(ctx.Request, http.MethodPost) {
			return
		}

		if !hasIdempotencyHeader(ctx.Request) {
			return
		}

		w := writter.NewWriter(ctx.Writer)
		key := keyFn(ctx.Request)

		defer func() {
			_ = cache.Unlock(key)
		}()

		if data, err := cache.Get(key); err != nil {
			if err := cache.Lock(key); isLocked(err) {
				ctx.Status(http.StatusConflict)
				ctx.Abort()
				return
			}
			ctx.Writer = w
			ctx.Next()
		} else {
			writeResponseAndAbort(ctx, data)
			return
		}

		if cacheCriteriaFn(ctx.Writer) {
			data := w.ToData(ctx.ContentType())
			_ = cache.Set(key, data)
		}
	}
}

func setOpts(opts ...Opts) (CacheKeyGenFn, CacheCriteriaFn) {
	keyFn := KeyAndURL
	cacheCriteriaFn := CacheOnlySuccess

	for _, opt := range opts {
		if opt.CacheKeyGenFn != nil {
			keyFn = opt.CacheKeyGenFn
		}

		if opt.CacheCriteriaFn != nil {
			cacheCriteriaFn = opt.CacheCriteriaFn
		}
	}

	return keyFn, cacheCriteriaFn
}

func isSomeHTTPMethod(r *http.Request, methods ...string) bool {
	for _, method := range methods {
		if r.Method == method {
			return true
		}
	}

	return false
}

func hasIdempotencyHeader(r *http.Request) bool {
	return r.Header.Get(idempotencyIDKey) != ""
}

func isLocked(err error) bool {
	return err != nil && errors.Is(err, cache.ErrAlreadyLocked)
}

func writeResponseAndAbort(ctx *gin.Context, data cache.Data) {
	data.WriteHeaders(ctx.Writer)
	ctx.Data(data.StatusCode, data.ContentType, data.Body)
	ctx.Abort()
}
