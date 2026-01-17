package views

import (
	"time"
)

type BaseView struct {
	UserName string
}

type ArticleShortItem struct {
	ID        int
	Title     string
	UpdatedAt time.Time
}

type ArticleFullViewItem struct {
	ID        int
	Title     string
	Content   string
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DashboardView struct {
	BaseView
	Articles []ArticleShortItem
}

type ArticleView struct {
	BaseView
	Article   ArticleFullViewItem
	CreatedBy string
}

type SearchView struct {
	DashboardView
	Query string
}
