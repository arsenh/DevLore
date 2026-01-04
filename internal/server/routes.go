package server

import (
	"fmt"
	"net/http"

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
	return routes
}

func (r *Routes) dashboardHandler(writer http.ResponseWriter, request *http.Request) {
	dashboardData, err := r.ArticleService.GetDashboardData(request.Context())
	if err != nil {
		panic("Error not handled yet.")
		//TODO: add correct error handling, maybe 404 page or something like that.
	}
	templates.RenderHTML(writer, templates.Dashboard, dashboardData)
}

func (r *Routes) rootHandler(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/dashboard", http.StatusPermanentRedirect)
}

func (r *Routes) newArticleHandler(writer http.ResponseWriter, request *http.Request) {
	templates.RenderHTML(writer, templates.NewArticle, nil)
}
