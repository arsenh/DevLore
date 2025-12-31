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

func DummyDatabase() []Article {
	now := time.Now()

	return []Article{
		{
			ID:        1,
			Title:     "Getting Started with Go",
			Content:   "Notes about installing Go, setting GOPATH, and writing the first program.",
			UserId:    11,
			CreatedAt: now.AddDate(0, 0, -10),
			UpdatedAt: now.AddDate(0, 0, -10),
		},
		{
			ID:        2,
			Title:     "Understanding Pointers in Go",
			Content:   "Pointers let you share and modify data without copying. & and * operators explained.",
			UserId:    22,
			CreatedAt: now.AddDate(0, 0, -9),
			UpdatedAt: now.AddDate(0, 0, -9),
		},
		{
			ID:        3,
			Title:     "HTTP Server Basics",
			Content:   "Minimal net/http example and how handlers work.",
			UserId:    33,
			CreatedAt: now.AddDate(0, 0, -8),
			UpdatedAt: now.AddDate(0, 0, -8),
		},
		{
			ID:        4,
			Title:     "Working with Structs",
			Content:   "Structs group related data together and are the building blocks of Go programs.",
			UserId:    44,
			CreatedAt: now.AddDate(0, 0, -7),
			UpdatedAt: now.AddDate(0, 0, -7),
		},
		{
			ID:        5,
			Title:     "Interfaces — Duck Typing in Go",
			Content:   "Interfaces describe behavior, not data. Any type that implements the methods satisfies it.",
			UserId:    55,
			CreatedAt: now.AddDate(0, 0, -6),
			UpdatedAt: now.AddDate(0, 0, -6),
		},
		{
			ID:        6,
			Title:     "Basic SQL with Go",
			Content:   "Connecting to a database, preparing statements, and scanning rows.",
			UserId:    66,
			CreatedAt: now.AddDate(0, 0, -5),
			UpdatedAt: now.AddDate(0, 0, -5),
		},
		{
			ID:        7,
			Title:     "Error Handling Patterns",
			Content:   "Why Go uses explicit errors and common patterns like wrapping and sentinel errors.",
			UserId:    77,
			CreatedAt: now.AddDate(0, 0, -4),
			UpdatedAt: now.AddDate(0, 0, -4),
		},
		{
			ID:        8,
			Title:     "Goroutines and Channels",
			Content:   "Concurrency basics: launching goroutines and communicating safely with channels.",
			UserId:    88,
			CreatedAt: now.AddDate(0, 0, -3),
			UpdatedAt: now.AddDate(0, 0, -3),
		},
		{
			ID:        9,
			Title:     "Template Rendering in Go",
			Content:   "Using html/template, parsing templates once, and passing data structs.",
			UserId:    99,
			CreatedAt: now.AddDate(0, 0, -2),
			UpdatedAt: now.AddDate(0, 0, -2),
		},
		{
			ID:        10,
			Title:     "Project Structure Best Practices",
			Content:   "Separating cmd/, internal/, handlers/, models/, templates/, and static/ folders.",
			UserId:    111,
			CreatedAt: now.AddDate(0, 0, -1),
			UpdatedAt: now.AddDate(0, 0, -1),
		},
	}
}
