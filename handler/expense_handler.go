package handler

import (
	"database/sql"
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
	category := getQueryUUID("category", c)
	dateFrom := getQueryDate("date-from", c)
	dateTo := getQueryDate("date-to", c)

	userIDValue, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "invalid user id"})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "invalid user id"})
		return
	}

	expenses, err := h.service.GetExpenses(limit, offset, category, dateFrom, dateTo, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with database"})
		return
	}
	total, err := h.service.GetTotalExpense(category, dateFrom, dateTo, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with database"})
		return
	}

	var res ExpensesResponse
	res.Expenses = make([]ExpenseResponse, 0, len(expenses))
	res.Limit = limit
	res.Offset = offset
	res.Total = total

	for _, expense := range expenses {
		res.Expenses = append(res.Expenses, ExpenseResponse{
			ID:          expense.ID,
			Amount:      expense.Amount,
			Note:        expense.Note,
			ExpenseDate: expense.ExpenseDate.Format(time.RFC3339),
			Category: Category{
				ID:           expense.Category.ID,
				CategoryName: expense.Category.CategoryName,
			},
		})
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

	userIDValue, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "invalid user id"})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "invalid user id"})
		return
	}

	expense := model.Expense{
		Amount:      req.Amount,
		Category:    model.Category{ID: *req.CategoryID},
		Note:        req.Note,
		ExpenseDate: expenseDate,
		User:        model.User{ID: userID},
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

func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	id := getParamUUID("id", c)
	if id == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid expense id"})
		return
	}

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

	expense := model.Expense{
		ID:          id,
		Amount:      req.Amount,
		Category:    model.Category{ID: *req.CategoryID},
		Note:        req.Note,
		ExpenseDate: expenseDate,
	}

	updatedExpense, err := h.service.UpdateExpense(expense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with database"})
		return
	}

	if updatedExpense == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "expense id not found"})
		return
	}

	res := ExpenseResponse{
		ID:          updatedExpense.ID,
		Amount:      updatedExpense.Amount,
		Note:        updatedExpense.Note,
		ExpenseDate: updatedExpense.ExpenseDate.Format(time.RFC3339),
		Category:    Category{ID: updatedExpense.Category.ID, CategoryName: updatedExpense.Category.CategoryName},
	}

	c.JSON(http.StatusOK, res)
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	id := getParamUUID("id", c)
	if id == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid expense id"})
		return
	}

	err := h.service.DeleteExpense(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "expense not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with database"})
		return
	}

	c.Status(http.StatusNoContent)
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

func getQueryUUID(name string, c *gin.Context) uuid.UUID {
	param := c.Query(name)
	return convertUUID(param)
}

func getParamUUID(name string, c *gin.Context) uuid.UUID {
	param := c.Param(name)
	return convertUUID(param)
}

func convertUUID(idStr string) uuid.UUID {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil
	}

	return id
}

func getQueryDate(name string, c *gin.Context) time.Time {
	param := c.Query(name)

	if param == "" {
		return time.Time{}
	}

	dateFormat := "2006-01-02"
	parsedDate, err := time.Parse(dateFormat, param)
	if err != nil {
		return time.Time{}
	}

	return parsedDate
}
