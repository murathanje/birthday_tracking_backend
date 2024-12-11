package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/murathanje/birthday_tracking_backend/internal/middleware"
	"github.com/murathanje/birthday_tracking_backend/internal/models"
	"github.com/murathanje/birthday_tracking_backend/internal/service"
)

type BirthdayHandler struct {
	birthdayService *service.BirthdayService
	categoryService *service.CategoryService
	userService     *service.UserService
}

func NewBirthdayHandler(birthdayService *service.BirthdayService, categoryService *service.CategoryService, userService *service.UserService) *BirthdayHandler {
	return &BirthdayHandler{
		birthdayService: birthdayService,
		categoryService: categoryService,
		userService:     userService,
	}
}

func (h *BirthdayHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	birthdays := api.Group("/birthdays")
	birthdays.Use(middleware.JWTAuth(func() []byte {
		return h.userService.GetJWTSecret()
	}))
	{
		birthdays.POST("", h.CreateBirthday)
		birthdays.GET("", h.GetUserBirthdays)
		birthdays.GET("/:id", h.GetBirthdayByID)
		birthdays.PUT("/:id", h.UpdateBirthday)
		birthdays.DELETE("/:id", h.DeleteBirthday)
	}
}

// CreateBirthday godoc
// @Summary Create a new birthday
// @Description Create a new birthday record for the authenticated user
// @Tags birthdays
// @Accept json
// @Produce json
// @Security Bearer
// @Param birthday body models.CreateBirthdayRequest true "Birthday details"
// @Success 201 {object} models.BirthdayResponse
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Server error"
// @Router /birthdays [post]
func (h *BirthdayHandler) CreateBirthday(c *gin.Context) {
	var req models.CreateBirthdayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Get user ID from JWT context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Validate category exists
	_, err = h.categoryService.GetCategoryByID(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	birthday, err := h.birthdayService.CreateBirthday(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, birthday.ToResponse())
}

// GetUserBirthdays godoc
// @Summary Get user's birthdays
// @Description Get all birthdays for the authenticated user
// @Tags birthdays
// @Produce json
// @Security Bearer
// @Success 200 {array} models.BirthdayResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Server error"
// @Router /birthdays [get]
func (h *BirthdayHandler) GetUserBirthdays(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	birthdays, err := h.birthdayService.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch birthdays"})
		return
	}

	response := make([]*models.BirthdayResponse, len(birthdays))
	for i, birthday := range birthdays {
		response[i] = birthday.ToResponse()
	}

	c.JSON(http.StatusOK, response)
}

// GetBirthdayByID godoc
// @Summary Get a birthday by ID
// @Description Get a birthday record by its ID (must belong to authenticated user)
// @Tags birthdays
// @Produce json
// @Security Bearer
// @Param id path string true "Birthday ID"
// @Success 200 {object} models.BirthdayResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not found"
// @Router /birthdays/{id} [get]
func (h *BirthdayHandler) GetBirthdayByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birthday ID"})
		return
	}

	birthday, err := h.birthdayService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Birthday not found"})
		return
	}

	// Verify ownership
	userID, _ := middleware.GetUserID(c)
	if birthday.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, birthday.ToResponse())
}

// UpdateBirthday godoc
// @Summary Update a birthday
// @Description Update a birthday record (must belong to authenticated user)
// @Tags birthdays
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Birthday ID"
// @Param birthday body models.CreateBirthdayRequest true "Birthday details"
// @Success 200 {object} models.BirthdayResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not found"
// @Router /birthdays/{id} [put]
func (h *BirthdayHandler) UpdateBirthday(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birthday ID"})
		return
	}

	var req models.CreateBirthdayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get existing birthday and verify ownership
	birthday, err := h.birthdayService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Birthday not found"})
		return
	}

	userID, _ := middleware.GetUserID(c)
	if birthday.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update fields
	parts := strings.Split(req.BirthDate, "-")
	if len(parts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birth date format"})
		return
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
		return
	}

	day, err := strconv.Atoi(parts[1])
	if err != nil || day < 1 || day > 31 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid day"})
		return
	}

	birthday.Name = req.Name
	birthday.BirthMonth = month
	birthday.BirthDay = day
	birthday.CategoryID = req.CategoryID
	birthday.Notes = req.Notes

	if err := h.birthdayService.Update(birthday); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update birthday"})
		return
	}

	c.JSON(http.StatusOK, birthday.ToResponse())
}

// DeleteBirthday godoc
// @Summary Delete a birthday
// @Description Delete a birthday record (must belong to authenticated user)
// @Tags birthdays
// @Produce json
// @Security Bearer
// @Param id path string true "Birthday ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not found"
// @Router /birthdays/{id} [delete]
func (h *BirthdayHandler) DeleteBirthday(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birthday ID"})
		return
	}

	// Get existing birthday and verify ownership
	birthday, err := h.birthdayService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Birthday not found"})
		return
	}

	userID, _ := middleware.GetUserID(c)
	if birthday.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.birthdayService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete birthday"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Birthday deleted successfully"})
}
