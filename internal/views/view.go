package views

type BaseView struct {
	UserName string
}

type ArticleShortItem struct {
	ID        string
	Title     string
	UpdatedAt string
}

type ArticleFullViewItem struct {
	ID        string
	Title     string
	Content   string
	UserId    string
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
