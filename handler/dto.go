package handler

import "github.com/google/uuid"

type ExpenseRequest struct {
	Amount      int       `json:"amount"`
	CategoryID  uuid.UUID `json:"category_id"`
	Note        string    `json:"note"`
	ExpenseDate string    `json:"expense_date"`
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
