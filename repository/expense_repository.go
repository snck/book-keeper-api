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

func (r *ExpenseRepository) GetExpenses(limit int, offset int) ([]model.Expense, error) {
	query := `
		SELECT
			e.id,
			e.amount,
			e.note,
			e.expense_date,
			e.created_at,
			e.updated_at,
			e.user_id,
			c.id,
			c.category_name
		FROM expenses e
		JOIN categories c ON c.id = e.category_id
		ORDER BY e.created_at DESC, e.id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make([]model.Expense, 0)
	for rows.Next() {
		var expense model.Expense
		var updatedAt sql.NullTime

		err := rows.Scan(
			&expense.ID,
			&expense.Amount,
			&expense.Note,
			&expense.ExpenseDate,
			&expense.CreatedAt,
			&updatedAt,
			&expense.User.ID,
			&expense.Category.ID,
			&expense.Category.CategoryName,
		)
		if err != nil {
			return nil, err
		}

		if updatedAt.Valid {
			expense.UpdatedAt = &updatedAt.Time
		}

		expenses = append(expenses, expense)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}
