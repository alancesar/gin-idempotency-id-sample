package writter

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"idempotency/internal/cache"
)

type (
	Writer struct {
		gin.ResponseWriter
		Body *bytes.Buffer
	}
)

func (w Writer) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w Writer) ToData(contentType string) cache.Data {
	return cache.Data{
		StatusCode:  w.Status(),
		ContentType: contentType,
		Headers:     w.Header().Clone(),
		Body:        w.Body.Bytes(),
	}
}

func NewWriter(w gin.ResponseWriter) *Writer {
	return &Writer{
		Body:           new(bytes.Buffer),
		ResponseWriter: w,
	}
}
