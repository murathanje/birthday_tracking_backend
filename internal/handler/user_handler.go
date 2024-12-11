package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/murathanje/birthday_tracking_backend/internal/config"
	"github.com/murathanje/birthday_tracking_backend/internal/middleware"
	"github.com/murathanje/birthday_tracking_backend/internal/models"
	"github.com/murathanje/birthday_tracking_backend/internal/service"
)

type UserHandler struct {
	service *service.UserService
	config  *config.Config
}

func NewUserHandler(service *service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		service: service,
		config:  cfg,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// Public routes (no authentication required)
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)

		// User routes (require JWT authentication)
		users := api.Group("/users")
		users.Use(middleware.JWTAuth(func() []byte {
			return []byte(h.config.JWTSecret)
		}))
		{
			users.GET("/me", h.GetCurrentUser)         // Get own profile
			users.PUT("/me", h.UpdateCurrentUser)      // Update own profile
			users.DELETE("/me", h.DeleteCurrentUser)   // Delete own account
		}

		// Admin routes (require API Key)
		admin := api.Group("/admin")
		admin.Use(middleware.APIKeyAuth(h.config))
		{
			admin.GET("/users", h.GetAllUsers)              // List all users
			admin.GET("/users/:id", h.GetUserByID)         // Get any user
			admin.PUT("/users/:id", h.UpdateUser)          // Update any user
			admin.DELETE("/users/:id", h.DeleteUser)       // Delete any user
		}
	}
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token for accessing protected endpoints
// @Description The returned token should be included in the Authorization header as "Bearer <token>"
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse "Successfully authenticated with JWT token"
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 401 {object} map[string]string "Invalid email or password"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	response, err := h.service.Login(&req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrInvalidCredentials {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account with email and password
// @Description After registration, use the /login endpoint to obtain a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User registration details"
// @Success 201 {object} models.UserResponse "User successfully registered"
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 409 {object} map[string]string "Email already exists"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Check if email already exists
	existingUser, err := h.service.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user.ToResponse())
}

// GetCurrentUser godoc
// @Summary Get current user's profile
// @Description Get the profile of the currently authenticated user
// @Tags users
// @Produce json
// @Security Bearer
// @Success 200 {object} models.UserResponse "User profile"
// @Failure 401 {object} map[string]string "Missing or invalid JWT token"
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateCurrentUser godoc
// @Summary Update current user's profile
// @Description Update the profile of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body models.UpdateUserRequest true "User update details"
// @Success 200 {object} models.UserResponse "Updated profile"
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 401 {object} map[string]string "Missing or invalid JWT token"
// @Router /users/me [put]
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := h.service.UpdateUser(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// DeleteCurrentUser godoc
// @Summary Delete current user's account
// @Description Delete the account of the currently authenticated user
// @Tags users
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]string "Success message"
// @Failure 401 {object} map[string]string "Missing or invalid JWT token"
// @Router /users/me [delete]
func (h *UserHandler) DeleteCurrentUser(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := h.service.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User account deleted successfully"})
}

// GetAllUsers godoc
// @Summary List all users
// @Description Get a list of all users (requires API Key)
// @Tags admin
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.UserResponse "List of users"
// @Failure 401 {object} map[string]string "Missing or invalid API Key"
// @Router /admin/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	response := make([]*models.UserResponse, len(users))
	for i, user := range users {
		response[i] = user.ToResponse()
	}

	c.JSON(http.StatusOK, response)
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Get a user by their ID (requires API Key)
// @Tags admin
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} models.UserResponse "User details"
// @Failure 401 {object} map[string]string "Missing or invalid API Key"
// @Failure 404 {object} map[string]string "User not found"
// @Router /admin/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateUser godoc
// @Summary Update any user
// @Description Update any user's information (requires API Key)
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID (UUID)"
// @Param user body models.UpdateUserRequest true "User update details"
// @Success 200 {object} models.UserResponse "Updated user"
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 401 {object} map[string]string "Missing or invalid API Key"
// @Failure 404 {object} map[string]string "User not found"
// @Router /admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := h.service.UpdateUser(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// DeleteUser godoc
// @Summary Delete any user
// @Description Delete any user account (requires API Key)
// @Tags admin
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} map[string]string "Success message"
// @Failure 401 {object} map[string]string "Missing or invalid API Key"
// @Failure 404 {object} map[string]string "User not found"
// @Router /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
} 