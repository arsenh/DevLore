package model

import "time"

type User struct {
	ID        int
	FullName  string
	Email     string
	Password  string
	CreatedAt time.Time
}
