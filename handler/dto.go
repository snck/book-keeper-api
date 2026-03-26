package handler

import "github.com/google/uuid"

type ExpenseRequest struct {
	Amount      int        `json:"amount" binding:"required"`
	CategoryID  *uuid.UUID `json:"category_id" binding:"required"`
	Note        string     `json:"note"`
	ExpenseDate string     `json:"expense_date" binding:"required"`
}

type ExpenseResponse struct {
	ID          uuid.UUID `json:"id"`
	Amount      int       `json:"amount"`
	Category    Category  `json:"category"`
	Note        string    `json:"note"`
	ExpenseDate string    `json:"expense_date"`
}

type Category struct {
	ID           uuid.UUID `json:"id"`
	CategoryName string    `json:"category_name"`
}

type ExpensesResponse struct {
	Expenses []ExpenseResponse `json:"expenses"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
	Total    int64             `json:"total"`
}

type SignupRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignupResponse struct {
	ID       uuid.UUID `json:"id"`
	UserName string    `json:"user_name"`
}

type LoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
	Token    string    `json:"token"`
}
