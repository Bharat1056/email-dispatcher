package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/user/queue/internal/context"
)

// Logger returns a gin.HandlerFunc that logs method, path, status code, and latency.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		ctx := context.NewContext(c)
		ctx.Next()

		latency := time.Since(start)
		log.Printf("%s %s %d %v",
			ctx.Method(),
			ctx.Path(),
			ctx.StatusCode(),
			latency,
		)
	}
}
