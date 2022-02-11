package writter

import (
	"bytes"
	"github.com/gin-gonic/gin"
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

func NewWriter(w gin.ResponseWriter) *Writer {
	return &Writer{
		Body:           new(bytes.Buffer),
		ResponseWriter: w,
	}
}
