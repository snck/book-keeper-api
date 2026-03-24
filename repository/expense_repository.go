package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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

func (r *ExpenseRepository) GetExpenses(limit int, offset int, categoryID uuid.UUID, dateFrom time.Time, dateTo time.Time) ([]model.Expense, error) {
	baseQuery := `
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
	`

	query, args := buildExpenseFilterQuery(baseQuery, categoryID, dateFrom, dateTo)

	args = append(args, limit, offset)
	query += fmt.Sprintf("\nORDER BY e.created_at DESC, e.id DESC\nLIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.Query(query, args...)
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

func (r *ExpenseRepository) GetTotalExpense(categoryID uuid.UUID, dateFrom time.Time, dateTo time.Time) (int64, error) {
	baseQuery := `
		SELECT COUNT(*)
		FROM expenses e
	`

	query, args := buildExpenseFilterQuery(baseQuery, categoryID, dateFrom, dateTo)

	var total int64
	err := r.db.QueryRow(query, args...).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *ExpenseRepository) UpdateExpense(expense model.Expense) (*model.Expense, error) {
	query := `
		WITH updated_expense AS (
			UPDATE expenses
			SET amount = $1, category_id = $2, note = $3, expense_date = $4, updated_at = NOW()
			WHERE id = $5
			RETURNING id, category_id
		)
		SELECT e.id, c.category_name
		FROM updated_expense e
		JOIN categories c ON c.id = e.category_id
	`

	err := r.db.QueryRow(
		query,
		expense.Amount,
		expense.Category.ID,
		expense.Note,
		expense.ExpenseDate,
		expense.ID,
	).Scan(&expense.ID, &expense.Category.CategoryName)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &expense, nil
}

func buildExpenseFilterQuery(baseQuery string, categoryID uuid.UUID, dateFrom time.Time, dateTo time.Time) (string, []any) {
	conditions := make([]string, 0, 3)
	args := make([]any, 0, 3)

	if categoryID != uuid.Nil {
		args = append(args, categoryID)
		conditions = append(conditions, fmt.Sprintf("e.category_id = $%d", len(args)))
	}

	if !dateFrom.IsZero() {
		args = append(args, dateFrom)
		conditions = append(conditions, fmt.Sprintf("e.expense_date >= $%d", len(args)))
	}

	if !dateTo.IsZero() {
		args = append(args, dateTo)
		conditions = append(conditions, fmt.Sprintf("e.expense_date < $%d", len(args)))
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += "\nWHERE " + strings.Join(conditions, "\n  AND ")
	}

	return query, args
}
