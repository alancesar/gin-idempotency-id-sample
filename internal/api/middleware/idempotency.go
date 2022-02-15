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

	KeyFn func(r *http.Request) interface{}
)

func DefaultKeyFn(r *http.Request) interface{} {
	return struct {
		Key string
		URL string
	}{
		Key: r.Header.Get(idempotencyIDKey),
		URL: r.RequestURI,
	}
}

func Idempotency(handler Cache, intentFn KeyFn) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !isSomeHTTPMethod(ctx.Request, http.MethodPost) {
			return
		}

		if !hasIdempotencyHeader(ctx.Request) {
			return
		}

		w := writter.NewWriter(ctx.Writer)
		key := intentFn(ctx.Request)

		defer func() {
			_ = handler.Unlock(key)
		}()

		if data, err := handler.Get(key); err != nil {
			if err := handler.Lock(key); isLocked(err) {
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

		if isSuccessStatusCode(ctx.Writer.Status()) {
			data := w.ToData(ctx.ContentType())
			_ = handler.Set(key, data)
		}
	}
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

func isSuccessStatusCode(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusBadRequest
}
