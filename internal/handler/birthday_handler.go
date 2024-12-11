package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/murathanje/birthday_tracking_backend/internal/models"
	"github.com/murathanje/birthday_tracking_backend/internal/service"
)

type BirthdayHandler struct {
	service *service.BirthdayService
}

func NewBirthdayHandler(service *service.BirthdayService) *BirthdayHandler {
	return &BirthdayHandler{service: service}
}

func (h *BirthdayHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		birthdays := api.Group("/birthdays")
		{
			birthdays.POST("", h.CreateBirthday)
			birthdays.GET("", h.GetAllBirthdays)
			birthdays.GET("/:id", h.GetBirthdayByID)
			birthdays.PUT("/:id", h.UpdateBirthday)
			birthdays.DELETE("/:id", h.DeleteBirthday)
			birthdays.GET("/upcoming", h.GetUpcomingBirthdays)
		}
	}
}

// CreateBirthday godoc
// @Summary Create a new birthday
// @Description Create a new birthday record
// @Tags birthdays
// @Accept json
// @Produce json
// @Param birthday body models.Birthday true "Birthday object"
// @Success 201 {object} models.Birthday
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /birthdays [post]
func (h *BirthdayHandler) CreateBirthday(c *gin.Context) {
	var birthday models.Birthday
	if err := c.ShouldBindJSON(&birthday); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if birthday.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	if birthday.BirthDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Birth date is required"})
		return
	}

	if err := h.service.CreateBirthday(&birthday); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create birthday"})
		return
	}

	c.JSON(http.StatusCreated, birthday)
}

// GetAllBirthdays godoc
// @Summary Get all birthdays
// @Description Get a list of all birthday records
// @Tags birthdays
// @Produce json
// @Success 200 {array} models.Birthday
// @Failure 500 {object} map[string]string
// @Router /birthdays [get]
func (h *BirthdayHandler) GetAllBirthdays(c *gin.Context) {
	birthdays, err := h.service.GetAllBirthdays()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch birthdays"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  birthdays,
		"count": len(birthdays),
	})
}

// GetBirthdayByID godoc
// @Summary Get a birthday by ID
// @Description Get a birthday record by its ID
// @Tags birthdays
// @Produce json
// @Param id path int true "Birthday ID"
// @Success 200 {object} models.Birthday
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /birthdays/{id} [get]
func (h *BirthdayHandler) GetBirthdayByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	birthday, err := h.service.GetBirthdayByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Birthday not found"})
		return
	}

	c.JSON(http.StatusOK, birthday)
}

// UpdateBirthday godoc
// @Summary Update a birthday
// @Description Update a birthday record by its ID
// @Tags birthdays
// @Accept json
// @Produce json
// @Param id path int true "Birthday ID"
// @Param birthday body models.Birthday true "Birthday object"
// @Success 200 {object} models.Birthday
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /birthdays/{id} [put]
func (h *BirthdayHandler) UpdateBirthday(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var birthday models.Birthday
	if err := c.ShouldBindJSON(&birthday); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if birthday.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	if birthday.BirthDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Birth date is required"})
		return
	}

	birthday.ID = uint(id)
	if err := h.service.UpdateBirthday(&birthday); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update birthday"})
		return
	}

	c.JSON(http.StatusOK, birthday)
}

// DeleteBirthday godoc
// @Summary Delete a birthday
// @Description Delete a birthday record by its ID
// @Tags birthdays
// @Produce json
// @Param id path int true "Birthday ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /birthdays/{id} [delete]
func (h *BirthdayHandler) DeleteBirthday(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteBirthday(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete birthday"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Birthday deleted successfully"})
}

// GetUpcomingBirthdays godoc
// @Summary Get upcoming birthdays
// @Description Get a list of birthdays in the next 30 days
// @Tags birthdays
// @Produce json
// @Success 200 {array} models.Birthday
// @Failure 500 {object} map[string]string
// @Router /birthdays/upcoming [get]
func (h *BirthdayHandler) GetUpcomingBirthdays(c *gin.Context) {
	birthdays, err := h.service.GetAllBirthdays()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch birthdays"})
		return
	}

	var upcomingBirthdays []models.Birthday
	now := time.Now()
	thirtyDaysFromNow := now.AddDate(0, 0, 30)

	for _, birthday := range birthdays {
		nextBirthday := time.Date(now.Year(), birthday.BirthDate.Month(), birthday.BirthDate.Day(), 0, 0, 0, 0, time.UTC)
		if nextBirthday.Before(now) {
			nextBirthday = nextBirthday.AddDate(1, 0, 0)
		}
		if nextBirthday.After(now) && nextBirthday.Before(thirtyDaysFromNow) {
			upcomingBirthdays = append(upcomingBirthdays, birthday)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  upcomingBirthdays,
		"count": len(upcomingBirthdays),
	})
}
