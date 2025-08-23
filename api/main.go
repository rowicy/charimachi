package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "template-mobile-app-api/docs"

	"github.com/joho/godotenv"

	util "template-mobile-app-api/util"
)

// @title Template Mobile App API
// @version 1.0
// @description This is a template API server using Go Gin framework
// @host localhost:8080
// @BasePath /api/v1

func main() {
	loadWarningIntersection()

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
		v1.GET("/directions/bicycle", util.GetDirections)
		// 目的地検索
		v1.GET("/search", util.GetSearch)
		//注意点
		v1.GET("/warningPoint", util.GetWarningPoints)
		//違反率
		v1.GET("/violation_rates", util.GetViolationRates)
	}

	// Start server on port 8080
	r.Run("0.0.0.0:8080")
}

func loadWarningIntersection() {
	warningIntersectionFile, err := os.Open("warningIntersection.json")
	if err != nil {
		fmt.Println("ファイルオープンエラー:", err)
		return
	}

	decoder := json.NewDecoder(warningIntersectionFile)
	if err := decoder.Decode(&util.WorningIntersectionPoints); err != nil {
		fmt.Println("JSONデコードエラー:", err)
		return
	}
	defer warningIntersectionFile.Close()
}
