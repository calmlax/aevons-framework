package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const DefaultHealthPath = "/health"

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service,omitempty"`
}

// RegisterHealthRoute registers the default health endpoint for Consul and probes.
func RegisterHealthRoute(r gin.IRoutes, serviceName string) {
	r.GET(DefaultHealthPath, func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthResponse{
			Status:  "ok",
			Service: serviceName,
		})
	})
}
