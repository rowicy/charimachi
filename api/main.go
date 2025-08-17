package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "template-mobile-app-api/docs"
)

// @title Template Mobile App API
// @version 1.0
// @description This is a template API server using Go Gin framework
// @host localhost:8080
// @BasePath /api/v1

// Post represents a post from JSONPlaceholder
type Post struct {
	UserID int    `json:"userId" example:"1"`
	ID     int    `json:"id" example:"1"`
	Title  string `json:"title" example:"Sample Post Title"`
	Body   string `json:"body" example:"Sample post body content"`
}

// PostResponse represents the formatted response
type PostResponse struct {
	ID          int    `json:"id" example:"1"`
	Title       string `json:"title" example:"Sample Post Title"`
	Body        string `json:"body" example:"Sample post body content"`
	UserID      int    `json:"userId" example:"1"`
	FormattedAt string `json:"formattedAt" example:"2024-01-01T12:00:00Z"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Internal server error"`
	Message string `json:"message" example:"Failed to fetch data"`
}

func main() {
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
		v1.GET("/posts", getPosts)
		v1.GET("/posts/:id", getPostByID)
		v1.GET("/health", getHealth)
	}

	// Start server on port 8080
	r.Run(":8080")
}

// getPosts godoc
// @Summary Get all posts from JSONPlaceholder
// @Description Fetch all posts from JSONPlaceholder API and return formatted response
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {array} PostResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts [get]
func getPosts(c *gin.Context) {
	// Fetch data from JSONPlaceholder
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to fetch data",
			Message: err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to read response",
			Message: err.Error(),
		})
		return
	}

	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to parse response",
			Message: err.Error(),
		})
		return
	}

	// Format the response
	var formattedPosts []PostResponse
	now := time.Now().Format(time.RFC3339)

	for _, post := range posts {
		formattedPosts = append(formattedPosts, PostResponse{
			ID:          post.ID,
			Title:       post.Title,
			Body:        post.Body,
			UserID:      post.UserID,
			FormattedAt: now,
		})
	}

	c.JSON(http.StatusOK, formattedPosts)
}

// getPostByID godoc
// @Summary Get a post by ID from JSONPlaceholder
// @Description Fetch a specific post by ID from JSONPlaceholder API and return formatted response
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} PostResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts/{id} [get]
func getPostByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid integer",
		})
		return
	}

	// Fetch data from JSONPlaceholder
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", id)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to fetch data",
			Message: err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Post not found",
			Message: fmt.Sprintf("Post with ID %d not found", id),
		})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to read response",
			Message: err.Error(),
		})
		return
	}

	var post Post
	if err := json.Unmarshal(body, &post); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to parse response",
			Message: err.Error(),
		})
		return
	}

	// Format the response
	formattedPost := PostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Body:        post.Body,
		UserID:      post.UserID,
		FormattedAt: time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, formattedPost)
}

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
