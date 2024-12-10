package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/murathanje/birthday_tracking_backend/internal/config"
	"github.com/murathanje/birthday_tracking_backend/internal/handler"
	"github.com/murathanje/birthday_tracking_backend/internal/middleware"
	"github.com/murathanje/birthday_tracking_backend/internal/repository"
	"github.com/murathanje/birthday_tracking_backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Server starting on %s in %s mode", serverAddr, os.Getenv("GIN_MODE"))
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
