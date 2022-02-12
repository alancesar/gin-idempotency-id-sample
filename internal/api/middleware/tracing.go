package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestIDKey = "X-Request-Id"
)

func Tracing(ctx *gin.Context) {
	if ctx.Request.Header.Get(requestIDKey) != "" {
		return
	}

	requestID := uuid.New().String()
	ctx.Header(requestIDKey, requestID)
}
