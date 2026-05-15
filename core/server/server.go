package server

import (
	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/xlog"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swagSpec "github.com/swaggo/swag"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service,omitempty"`
}

// RegisterHealthRoute registers the default health endpoint for Consul and probes.
func RegisterHealthRoute(r gin.IRoutes, serviceName string) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthResponse{
			Status:  "ok",
			Service: serviceName,
		})
	})
}

// RegisterOpenApiRoute registers the default OpenAPI endpoint for Consul and probes.
func RegisterOpenApiRoute(r gin.IRoutes, cfg *config.Config) {
	r.GET("/api/swagger.json", func(c *gin.Context) {
		doc, err := loadOpenAPIDoc()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(doc))
	})
	xlog.Info("swagger json:      http://localhost:%d/api/swagger.json", cfg.Server.Port)
}

func loadOpenAPIDoc() (string, error) {
	doc, err := swagSpec.ReadDoc()
	if err == nil && doc != "" && doc != "{}" {
		return doc, nil
	}

	data, readErr := os.ReadFile("api/swagger.json")
	if readErr != nil {
		if err != nil {
			return "", err
		}
		return "", readErr
	}

	return string(data), nil
}
