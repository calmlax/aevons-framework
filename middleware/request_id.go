package middleware

import (
	"crypto/rand"
	"fmt"

	"github.com/gin-gonic/gin"
)

const RequestIdKey = "X-Request-ID"

// newUUID generates a random UUID v4 string.
func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant bits
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// RequestID returns a Gin middleware that assigns a unique ID to each request,
// stores it in the context, and sets the X-Request-ID response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := newUUID()
		c.Set(RequestIdKey, id)
		c.Header(RequestIdKey, id)
		c.Next()
	}
}
