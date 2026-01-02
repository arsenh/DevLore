package model

import "time"

type Article struct {
	ID        int
	Title     string
	Content   string
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
