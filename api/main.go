package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "template-mobile-app-api/docs"

	"github.com/joho/godotenv"
)

// @title Template Mobile App API
// @version 1.0
// @description This is a template API server using Go Gin framework
// @host localhost:8080
// @BasePath /api/v1

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Internal server error"`
	Message string `json:"message" example:"Failed to fetch data"`
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", getHealth)
		// 経路検索
		v1.GET("/directions/bicycle", getDirections)
		// 目的地検索
		v1.GET("/search", getSearch)
	}

	// Start server on port 8080
	r.Run(":8080")
}
