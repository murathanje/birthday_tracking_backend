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

func (h *BirthdayHandler) GetAllBirthdays(c *gin.Context) {
	birthdays, err := h.service.GetAllBirthdays()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch birthdays"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": birthdays,
		"count": len(birthdays),
	})
}

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

	// Validate required fields
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

func (h *BirthdayHandler) GetUpcomingBirthdays(c *gin.Context) {
	// Get upcoming birthdays in the next 30 days
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
		"data": upcomingBirthdays,
		"count": len(upcomingBirthdays),
	})
} 