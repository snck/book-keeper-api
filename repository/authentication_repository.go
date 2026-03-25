package repository

import (
	"database/sql"

	"github.com/snck/book-keeper-api/model"
)

type AuthenticationRepository struct {
	db *sql.DB
}

func NewAuthenticationRepository(db *sql.DB) *AuthenticationRepository {
	return &AuthenticationRepository{db: db}
}

func (r *AuthenticationRepository) AddNewUser(user model.User) (*model.User, error) {
	query := `
		INSERT INTO users (user_name, password_hash)
		VALUES ($1, $2)
		ON CONFLICT (user_name) DO NOTHING
		RETURNING id, user_name
	`

	err := r.db.QueryRow(
		query,
		user.UserName,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.UserName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
