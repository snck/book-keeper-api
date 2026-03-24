package model

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID           uuid.UUID
	CategoryName string
	User         User
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}
