package server

import (
	"fmt"
	"net/http"

	"github.com/arsenh/DevLore/internal/config"
	"github.com/arsenh/DevLore/internal/templates"
)

func GetRoutes() *http.ServeMux {
	routes := http.NewServeMux()

	// setup file server
	files := http.FileServer(http.Dir(config.PublicDir))
	routes.Handle(fmt.Sprintf("/%s/", config.StaticDir), http.StripPrefix(fmt.Sprintf("/%s/", config.StaticDir), files))

	routes.HandleFunc("/", rootHandler)
	routes.HandleFunc("/dashboard", dashboardHandler)
	routes.HandleFunc("/articles/new", newArticleHandler)
	return routes
}

func dashboardHandler(writer http.ResponseWriter, request *http.Request) {
	templates.RenderHTML(writer, templates.Dashboard, nil)
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/dashboard", http.StatusPermanentRedirect)
}

func newArticleHandler(writer http.ResponseWriter, request *http.Request) {
	templates.RenderHTML(writer, templates.NewArticle, nil)
}
