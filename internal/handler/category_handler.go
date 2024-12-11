package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/murathanje/birthday_tracking_backend/internal/models"
	"github.com/murathanje/birthday_tracking_backend/internal/service"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		categories := api.Group("/categories")
		{
			categories.GET("", h.GetAllCategories)
			categories.GET("/:id", h.GetCategoryByID)
			categories.GET("/name/:name", h.GetCategoryByName)
		}
	}
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Get a list of all birthday categories with their icons
// @Description Each category includes an ID, name, and an emoji icon
// @Description Use this endpoint to populate category dropdowns or lists
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category
// @Failure 500 {object} map[string]string
// @Router /categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	// Use the models.Category type explicitly to satisfy the unused import
	var _ []models.Category = categories

	c.JSON(http.StatusOK, categories)
}

// GetCategoryByID godoc
// @Summary Get a category by ID
// @Description Get a category by its ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID (UUID)"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCategoryByName godoc
// @Summary Get a category by name
// @Description Get a category by its name (e.g., 'family', 'friends', etc.)
// @Tags categories
// @Produce json
// @Param name path string true "Category name"
// @Success 200 {object} models.Category
// @Failure 404 {object} map[string]string
// @Router /categories/name/{name} [get]
func (h *CategoryHandler) GetCategoryByName(c *gin.Context) {
	name := c.Param("name")
	category, err := h.service.GetCategoryByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
} 