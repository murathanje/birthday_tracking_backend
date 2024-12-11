package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/murathanje/birthday_tracking_backend/docs"
	"github.com/murathanje/birthday_tracking_backend/internal/config"
	"github.com/murathanje/birthday_tracking_backend/internal/handler"
	"github.com/murathanje/birthday_tracking_backend/internal/middleware"
	"github.com/murathanje/birthday_tracking_backend/internal/repository"
	"github.com/murathanje/birthday_tracking_backend/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Birthday Tracking API
// @version         1.0
// @description     A birthday tracking service API in Go using Gin framework.
// @description     Features:
// @description     - User management with JWT authentication for user operations
// @description     - API Key authentication for admin operations
// @description     - Birthday tracking with simple categorization (string-based)
// @description     - Example categories: "Family", "Friend", "Work", "School", etc.
// @description     - Upcoming birthdays tracking
// @description     
// @description     Authentication:
// @description     1. For Users:
// @description        - Register a new account using /api/v1/register
// @description        - Login with your credentials at /api/v1/login to get a JWT token
// @description        - Use the token in the Authorization header for protected endpoints
// @description        - Format: "Bearer <your_jwt_token>"
// @description     2. For Admins:
// @description        - Use API Key in the X-API-Key header for admin endpoints
// @description        - The API Key should be set in your .env file
// @description     
// @description     Endpoints:
// @description     1. Auth Endpoints (Public):
// @description        - POST /api/v1/register - Create new account
// @description        - POST /api/v1/login - Get JWT token
// @description     2. User Endpoints (Requires JWT):
// @description        - GET /api/v1/users/me - Get own profile
// @description        - PUT /api/v1/users/me - Update own profile
// @description        - DELETE /api/v1/users/me - Delete own account
// @description     3. Admin Endpoints (Requires API Key):
// @description        - GET /api/v1/admin/users - List all users
// @description        - GET /api/v1/admin/users/{id} - Get any user
// @description        - PUT /api/v1/admin/users/{id} - Update any user
// @description        - DELETE /api/v1/admin/users/{id} - Delete any user
// @description     4. Birthday Endpoints (Requires JWT):
// @description        - POST /api/v1/birthdays - Create birthday (with category as string)
// @description        - GET /api/v1/birthdays - List own birthdays
// @description        - GET /api/v1/birthdays/{id} - Get specific birthday
// @description        - PUT /api/v1/birthdays/{id} - Update birthday
// @description        - DELETE /api/v1/birthdays/{id} - Delete birthday
// @description     
// @description     Birthday Categories:
// @description     Categories are now implemented as simple strings. You can use any string value
// @description     for categorization. Some suggested categories:
// @description     - "Family" - For family members
// @description     - "Friend" - For friends
// @description     - "Work" - For work colleagues
// @description     - "School" - For school/university friends
// @description     - "Other" - For any other category

// @contact.name   API Support
// @contact.url    https://github.com/murathanje/birthday_tracking_backend
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:5050
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Required for user-specific operations.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API Key required for admin operations. Set this in your .env file.

// @tag.name auth
// @tag.description Authentication endpoints for user registration and login

// @tag.name users
// @tag.description User-specific endpoints (requires JWT authentication)

// @tag.name admin
// @tag.description Admin endpoints for user management (requires API Key)

// @tag.name birthdays
// @tag.description Birthday management endpoints with string-based categorization (requires JWT authentication)

// @schemes http https

func main() {
	gin.SetMode(gin.ReleaseMode)
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	birthdayRepo := repository.NewBirthdayRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, cfg)
	birthdayService := service.NewBirthdayService(birthdayRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, cfg)
	birthdayHandler := handler.NewBirthdayHandler(birthdayService, userService)

	router := gin.New()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Register routes
	userHandler.RegisterRoutes(router)
	birthdayHandler.RegisterRoutes(router)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Server starting on %s in %s mode", serverAddr, os.Getenv("GIN_MODE"))
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
