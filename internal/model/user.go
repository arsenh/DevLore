package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	FullName  string
	Email     string
	Password  string
	CreatedAt time.Time
}
