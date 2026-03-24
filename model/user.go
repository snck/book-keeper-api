package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	UserName  string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
