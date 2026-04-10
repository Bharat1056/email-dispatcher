package router

import (
	"github.com/gin-gonic/gin"

	"github.com/user/queue/internal/context"
	"github.com/user/queue/internal/handler"
	"github.com/user/queue/internal/middleware"
)

// Setup creates and configures the Gin engine with middleware and routes.
func Setup() *gin.Engine {
	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", context.Wrap(handler.HealthCheck))

	return r
}
