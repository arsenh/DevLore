package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/arsenh/DevLore/internal/config"
	"github.com/arsenh/DevLore/internal/service"
	"github.com/arsenh/DevLore/internal/templates"
)

type Routes struct {
	articleService *service.ArticleService
}

func NewRoutes(service *service.ArticleService) *Routes {
	return &Routes{
		articleService: service,
	}
}

func (r *Routes) GetRoutes() *http.ServeMux {
	routes := http.NewServeMux()

	// setup file server
	files := http.FileServer(http.Dir(config.PublicDir))
	routes.Handle(fmt.Sprintf("/%s/", config.StaticDir), http.StripPrefix(fmt.Sprintf("/%s/", config.StaticDir), files))

	routes.HandleFunc("/", r.rootHandler)
	routes.HandleFunc("/dashboard", r.dashboardHandler)
	routes.HandleFunc("/articles/new", r.newArticleHandler)
	routes.HandleFunc("/articles/", r.viewArticleHandler)
	return routes
}

func (r *Routes) dashboardHandler(writer http.ResponseWriter, request *http.Request) {
	dashboardData, err := r.articleService.GetDashboardData(request.Context())
	if err != nil {
		templates.InternalServerError(writer, err)
	}
	if err := templates.Render(writer, http.StatusOK, templates.DashboardTemplate, dashboardData); err != nil {
		templates.InternalServerError(writer, err)
	}
}

func (r *Routes) rootHandler(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/dashboard", http.StatusPermanentRedirect)
}

func (r *Routes) newArticleHandler(writer http.ResponseWriter, request *http.Request) {
	if (request.Method != http.MethodGet) && (request.Method != http.MethodPost) {
		templates.BadRequest(writer)
	}

	ctx := request.Context()

	switch request.Method {
	case http.MethodGet:
		// render form for new article creation
		if err := templates.Render(writer, http.StatusOK, templates.NewArticleTemplate, nil); err != nil {
			templates.InternalServerError(writer, err)
		}
		return
	case http.MethodPost:
		// parse user input data to create new article
		if err := request.ParseForm(); err != nil {
			templates.BadRequest(writer)
			return
		}
		//TODO: need to get also UserId which created the article.
		title := request.Form.Get("title")
		content := request.Form.Get("content")
		id, err := r.articleService.SaveArticle(ctx, title, content)
		if err != nil {
			templates.InternalServerError(writer, err)
			return
		}
		view, err := r.articleService.GetArticleById(ctx, id)
		if err != nil {
			//this case is internal error, because, we must get article that created previously
			templates.InternalServerError(writer, err)
			return
		}
		err = templates.Render(writer, http.StatusOK, templates.ViewArticleTemplate, view)
		return
	default:
		templates.BadRequest(writer)
	}

}

func (r *Routes) viewArticleHandler(writer http.ResponseWriter, request *http.Request) {
	path := strings.TrimPrefix(request.URL.Path, "/articles/")
	idStr := strings.Trim(path, "/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		templates.BadRequest(writer)
		return
	}

	data, err := r.articleService.GetArticleById(request.Context(), id)
	if err != nil {
		templates.NotFound(writer)
		return
	}

	if err := templates.Render(writer, http.StatusOK, templates.ViewArticleTemplate, data); err != nil {
		templates.InternalServerError(writer, err)
	}
}
