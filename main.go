package main

import (
	"net/http"

	"github.com/KingLeak95/todo-list-go/docs"
	"github.com/KingLeak95/todo-list-go/middleware"
	"github.com/KingLeak95/todo-list-go/models"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Todo List API
// @version 1.0
// @description A modern Todo List API with comprehensive features
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	r := gin.Default()

	// Enable Middle Ware from Gin Package
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Security middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RateLimitMiddleware(10, 100)) // 10 requests per second, burst of 100

	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "Todo List API"
	docs.SwaggerInfo.Description = "A modern Todo List API with comprehensive features"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Connect Database
	models.ConnectDatabase()

	// Index for Testing
	// @Summary Health check endpoint
	// @Description Returns a simple health check response
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router / [get]
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Todo List API is running", "version": "1.0"})
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public authentication routes
	r.POST("/auth/register", models.Register)
	r.POST("/auth/login", models.Login)
	r.POST("/auth/refresh", models.RefreshToken)

	// Legacy user creation (deprecated, use /auth/register)
	r.POST("/createUser", models.CreateUser)

	// Protected routes - require authentication
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Users
		protected.GET("/allUsers", models.GetAllUsers)
		protected.DELETE("/deleteUser/:id", models.DeleteUser)

		// Tasks
		protected.GET("/tasks", models.GetAllTasks)
		protected.POST("/tasks", models.CreateTask)
		protected.GET("/tasks/:id", models.GetTaskByID)
		protected.PUT("/tasks/:id", models.UpdateTask)
		protected.DELETE("/tasks/:id", models.DeleteTask)
		protected.PUT("/tasks/:id/complete", models.CompleteTask)
		protected.GET("/users/:id/tasks", models.GetTasksByUser)
	}

	r.Run()
}
