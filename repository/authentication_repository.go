package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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

func (r *AuthenticationRepository) GetUserByUserName(userName string) (*model.User, error) {
	query := `
		SELECT id, user_name, password_hash, created_at, updated_at
		FROM users
		WHERE user_name = $1
	`

	var user model.User
	err := r.db.QueryRow(query, userName).Scan(
		&user.ID,
		&user.UserName,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthenticationRepository) AddTokenToBlocklist(userID uuid.UUID, token string, expiredAt time.Time) error {
	query := `
		INSERT INTO blocklists (user_id, token, expired_at)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(query, userID, token, expiredAt)
	return err
}
