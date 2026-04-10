package handler

import (
	"github.com/user/queue/internal/context"
)

// HealthCheck returns a simple JSON response indicating the service is running.
func HealthCheck(c *context.Context) {
	c.OK(map[string]string{"status": "ok"})
}
