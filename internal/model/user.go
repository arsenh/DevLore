package model

import "time"

type User struct {
	ID        int
	FullName  string
	Email     string
	Password  string
	CreatedAt time.Time
}

func DummyUsers() []User {
	return []User{
		{ID: 11, FullName: "Alice Johnson"},
		{ID: 22, FullName: "Bob Smith"},
		{ID: 33, FullName: "Carol Adams"},
		{ID: 44, FullName: "Dave Brown"},
		{ID: 55, FullName: "Eve Thompson"},
		{ID: 66, FullName: "Frank Williams"},
		{ID: 77, FullName: "Grace Miller"},
		{ID: 88, FullName: "Heidi Clark"},
		{ID: 99, FullName: "Ivan Garcia"},
		{ID: 111, FullName: "Judy Martinez"},
	}
}
