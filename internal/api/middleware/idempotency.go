package middleware

import (
	"github.com/gin-gonic/gin"
	"idempotency/internal/api/writter"
	"idempotency/internal/cache"
	"net/http"
)

type (
	Cache interface {
		Get(idempotencyID, url string) (cache.Data, error)
		Set(idempotencyID, url string, data cache.Data) error
	}
)

func Idempotency(handler Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodPost {
			return
		}

		idempotencyID := ctx.Request.Header.Get("Idempotency-Id")
		if idempotencyID == "" {
			return
		}

		w := writter.NewWriter(ctx.Writer)
		if data, err := handler.Get(idempotencyID, ctx.Request.RequestURI); err != nil {
			ctx.Writer = w
			ctx.Next()
		} else {
			writeHeaders(ctx.Writer, data)
			ctx.Data(data.StatusCode, data.ContentType, data.Body)
			ctx.Abort()
			return
		}

		statusCode := ctx.Writer.Status()
		if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
			_ = handler.Set(idempotencyID, ctx.Request.RequestURI,
				cache.Data{
					StatusCode:  statusCode,
					ContentType: ctx.ContentType(),
					Headers:     ctx.Writer.Header().Clone(),
					Body:        w.Body.Bytes(),
				})
		}
	}
}

func writeHeaders(writer http.ResponseWriter, data cache.Data) {
	for k, v := range data.Headers {
		writer.Header().Set(k, v[0])
	}
}
