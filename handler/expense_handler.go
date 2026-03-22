package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snck/book-keeper-api/service"
)

type ExpenseHandler struct {
	service *service.ExpenseService
}

func NewExpenseHandler(service *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
