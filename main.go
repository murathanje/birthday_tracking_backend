package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/murathanje/birthday_tracking_backend/internal/api"
	"github.com/murathanje/birthday_tracking_backend/internal/config"
	"github.com/murathanje/birthday_tracking_backend/internal/repository"
	"github.com/murathanje/birthday_tracking_backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Setup database connection
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repository, service, and handler
	birthdayRepo := repository.NewBirthdayRepository(db)
	birthdayService := service.NewBirthdayService(birthdayRepo)
	birthdayHandler := api.NewBirthdayHandler(birthdayService)

	// Initialize router
	router := gin.Default()

	// Register routes
	birthdayHandler.RegisterRoutes(router)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Start server
	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}