package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/snck/book-keeper-api/model"
	"github.com/snck/book-keeper-api/service"
)

type ExpenseHandler struct {
	service *service.ExpenseService
}

func NewExpenseHandler(service *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {

	limit := getQueryInt("limit", 10, c)
	offset := getQueryInt("offset", 0, c)

	res, err := h.service.GetExpenses(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with database"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var req ExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	expenseDate, err := time.Parse(time.RFC3339, req.ExpenseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid time format"})
		return
	}

	tempUserID, err := uuid.Parse("00000000-0000-0000-0000-000000000001")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error parsing id"})
		return
	}

	expense := model.Expense{
		Amount:      req.Amount,
		Category:    model.Category{ID: req.CategoryID},
		Note:        req.Note,
		ExpenseDate: expenseDate,
		User:        model.User{ID: tempUserID},
	}

	expense, err = h.service.CreateExpense(expense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with database"})
		return
	}

	res := ExpenseResponse{
		ID:          expense.ID,
		Amount:      expense.Amount,
		Note:        expense.Note,
		ExpenseDate: expense.ExpenseDate.Format(time.RFC3339),
		Category:    Category{ID: expense.Category.ID, CategoryName: expense.Category.CategoryName},
	}

	c.JSON(http.StatusCreated, res)
}

func getQueryInt(name string, defaultValue int, c *gin.Context) int {
	param := c.Query(name)

	if param == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}

	return parsedValue
}
