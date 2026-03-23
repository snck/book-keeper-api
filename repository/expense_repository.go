package repository

import (
	"database/sql"

	"github.com/snck/book-keeper-api/model"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) CreateExpense(expense model.Expense) (model.Expense, error) {
	query := `
		WITH inserted_expense AS (
			INSERT INTO expenses (amount, category_id, note, user_id, expense_date)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, category_id
		)
		SELECT e.id, c.category_name
		FROM inserted_expense e
		JOIN categories c ON c.id = e.category_id
	`

	err := r.db.QueryRow(
		query,
		expense.Amount,
		expense.Category.ID,
		expense.Note,
		expense.User.ID,
		expense.ExpenseDate,
	).Scan(&expense.ID, &expense.Category.CategoryName)

	if err != nil {
		return expense, err
	}

	return expense, nil
}
