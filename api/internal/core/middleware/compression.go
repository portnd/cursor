package middleware

import (
	"compress/gzip"
	"strings"

	"github.com/gin-gonic/gin"
)

// GzipResponse compresses JSON/text responses when the client accepts gzip.
func GzipResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")
		gw := gzip.NewWriter(c.Writer)
		defer gw.Close()

		c.Writer = &gzipResponseWriter{ResponseWriter: c.Writer, Writer: gw}
		c.Next()
	}
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func (w *gzipResponseWriter) WriteString(s string) (int, error) {
	return w.Writer.Write([]byte(s))
}
