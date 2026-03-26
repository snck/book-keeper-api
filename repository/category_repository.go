package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/snck/book-keeper-api/model"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetCategories(userID uuid.UUID) ([]model.Category, error) {
	query := `
		SELECT id, category_name
		FROM categories
		WHERE user_id = $1
		ORDER BY category_name
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]model.Category, 0)
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.ID, &category.CategoryName)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
