package views

type BaseView struct {
	UserName string
}

type ArticleShortItem struct {
	ID        int
	Title     string
	UpdatedAt string
}

type ArticleFullViewItem struct {
	ID        int
	Title     string
	Content   string
	UserId    int
	CreatedAt string
	UpdatedAt string
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
