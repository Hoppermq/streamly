package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	domain "github.com/hoppermq/streamly/pkg/domain/core"
)

type HealthResponse struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
}

func HealthHandler(callback domain.HealthCallback) gin.HandlerFunc {
	return func(c *gin.Context) {
		isHealthy, err := callback(c.Request.Context())

		resp := HealthResponse{
			Data:      processHealthy(isHealthy),
			Timestamp: time.Now(),
		}
		var statusCode int

		if err != nil {
			resp.Error = err.Error()
			statusCode = http.StatusInternalServerError
		} else if !isHealthy {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, resp)
	}
}

func processHealthy(isHealthy bool) string {
	if isHealthy {
		return "ok"
	}

	return "unhealthy"
}
