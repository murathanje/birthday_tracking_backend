package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/murathanje/birthday_tracking_backend/docs" // Import swagger docs
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


// @host      localhost:5050
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	gin.SetMode(gin.ReleaseMode)
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	birthdayRepo := repository.NewBirthdayRepository(db)
	birthdayService := service.NewBirthdayService(birthdayRepo)
	birthdayHandler := handler.NewBirthdayHandler(birthdayService)

	router := gin.New()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	birthdayHandler.RegisterRoutes(router)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Server starting on %s in %s mode", serverAddr, os.Getenv("GIN_MODE"))
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
