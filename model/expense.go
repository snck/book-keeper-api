package model

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID          uuid.UUID
	Amount      int
	Category    Category
	Note        string
	User        User
	ExpenseDate time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
