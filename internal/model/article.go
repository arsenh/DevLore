package model

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID        uuid.UUID
	Title     string
	Content   string
	UserId    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
