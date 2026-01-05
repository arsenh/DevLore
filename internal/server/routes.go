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
	ArticleService *service.ArticleService
}

func NewRoutes(articleService *service.ArticleService) *Routes {
	return &Routes{
		ArticleService: articleService,
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
	dashboardData, err := r.ArticleService.GetDashboardData(request.Context())
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
	templates.Render(writer, http.StatusOK, templates.NewArticleTemplate, nil)
}

func (r *Routes) viewArticleHandler(writer http.ResponseWriter, request *http.Request) {
	path := strings.TrimPrefix(request.URL.Path, "/articles/")
	idStr := strings.Trim(path, "/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		templates.BadRequest(writer)
		return
	}

	data, err := r.ArticleService.GetArticleById(request.Context(), id)
	if err != nil {
		templates.NotFound(writer)
		return
	}

	if err := templates.Render(writer, http.StatusOK, templates.ViewArticleTemplate, data); err != nil {
		templates.InternalServerError(writer, err)
	}
}
