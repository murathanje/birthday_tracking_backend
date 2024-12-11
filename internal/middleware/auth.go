package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/murathanje/birthday_tracking_backend/internal/config"
)

// APIKeyAuth creates a middleware for system-level API key authentication
func APIKeyAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key") // Using X-API-Key header for API key
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "X-API-Key header is required"})
			c.Abort()
			return
		}

		if apiKey != cfg.APIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// JWTAuth creates a middleware for user-level JWT authentication
func JWTAuth(getSecret func() []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract the token from the Authorization header
		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &jwt.MapClaims{}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return getSecret(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract user information from claims
		userIDStr, ok := (*claims)["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		email, ok := (*claims)["email"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", userID)
		c.Set("user_email", email)

		c.Next()
	}
}

// GetUserID extracts the user ID from the context
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return userID.(uuid.UUID), nil
}

// RequireOwner creates a middleware for checking if the user is the owner of the resource
func RequireOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		resourceUserID := c.Param("id")
		if resourceUserID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			c.Abort()
			return
		}

		c.Next()
	}
} 