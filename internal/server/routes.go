package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arsenh/DevLore/internal/config"
	"github.com/arsenh/DevLore/internal/logger"
	"github.com/arsenh/DevLore/internal/service"
	"github.com/arsenh/DevLore/internal/templates"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Routes struct {
	articleService *service.ArticleService
}

func NewRoutes(service *service.ArticleService) *Routes {
	return &Routes{
		articleService: service,
	}
}

func (r *Routes) GetRoutes() http.Handler {
	router := chi.NewRouter()

	// setup logrus for requests
	logger := logger.L()
	chiLogger := middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger: logger,
	})

	router.Use(chiLogger)

	fs := http.FileServer(http.Dir(config.PublicDir))
	router.Handle(fmt.Sprintf("/%s/*", config.StaticDir), http.StripPrefix(fmt.Sprintf("/%s/", config.StaticDir), fs))

	// global 404 not found page
	router.NotFound(r.notFoundPage)

	router.Get("/", r.rootHandler)
	router.Get("/dashboard", r.dashboardHandler)
	router.Get("/articles/{id}", r.viewArticleHandler)
	router.Get("/articles/new", r.showNewArticleHandler)
	router.Post("/articles/new", r.createNewArticleHandler)
	return router
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

func (r *Routes) showNewArticleHandler(writer http.ResponseWriter, request *http.Request) {
	// render form for new article creation
	if err := templates.Render(writer, http.StatusOK, templates.NewArticleTemplate, nil); err != nil {
		templates.InternalServerError(writer, err)
	}
}

func (r *Routes) createNewArticleHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

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
	http.Redirect(writer, request, fmt.Sprintf("/articles/%d", id), http.StatusSeeOther)
}

func (r *Routes) viewArticleHandler(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")

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

func (r *Routes) notFoundPage(writer http.ResponseWriter, request *http.Request) {
	templates.NotFound(writer)
}
