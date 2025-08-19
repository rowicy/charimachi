package main
import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

// getHealth godoc
// @Summary Health check endpoint
// @Description Returns the health status of the API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "template-mobile-app-api",
		"version":   "1.0.0",
	})
}