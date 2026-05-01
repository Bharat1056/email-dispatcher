package handler

import (
	"github.com/user/queue/internal/context"
)

type HealthResponse struct {
	Status string `json:"status"`
}

// HealthCheck returns a simple JSON response indicating the service is running.
func HealthCheck(c *context.Context) {
	c.OK(HealthResponse{Status: "ok"})
}
