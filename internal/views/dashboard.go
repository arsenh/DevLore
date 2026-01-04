package views

import (
	"time"
)

type ArticleItem struct {
	ID        int
	Title     string
	UpdatedAt time.Time
}

type DashboardView struct {
	UserName string
	Articles []ArticleItem
}
