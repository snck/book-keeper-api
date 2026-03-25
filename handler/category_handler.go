package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snck/book-keeper-api/service"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with database"})
		return
	}

	response := make([]Category, 0, len(categories))
	for _, category := range categories {
		response = append(response, Category{
			ID:           category.ID,
			CategoryName: category.CategoryName,
		})
	}

	c.JSON(http.StatusOK, response)
}
