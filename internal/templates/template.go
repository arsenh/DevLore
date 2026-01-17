package templates

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/arsenh/DevLore/internal/logger"
)

//go:embed *.html
var templatesFS embed.FS

const (
	LayoutTemplate              = "layout"
	DashboardTemplate           = "dashboard.html"
	NewArticleTemplate          = "new-article.html"
	EditArticleTemplate         = "edit-article.html"
	ViewArticleTemplate         = "view-article.html"
	SearchTemplate              = "search.html"
	NotFoundTemplate            = "404.html"
	BadRequestTemplate          = "400.html"
	InternalServerErrorTemplate = "500.html"
	RegisterTemplate            = "register.html"
	LoginTemplate               = "login.html"
	NotPermitted                = "not-permitted.html"
)

var base *template.Template

func init() {
	var err error

	base, err = template.ParseFS(
		templatesFS,
		"layout.html",
	)

	if err != nil {
		panic(fmt.Errorf("parse base template: %w", err))
	}

	logger.L().Info("templates are embedded and parsed successfully")
}

func Render(w http.ResponseWriter, status int, page string, data interface{}) error {
	//Clone the base layout template
	tmpl, err := base.Clone()
	if err != nil {
		return fmt.Errorf("clone base template %q: %w", base.Name(), err)
	}

	//Parse the page template into the clone
	if _, err := tmpl.ParseFS(templatesFS, page); err != nil {
		return fmt.Errorf("parse page template %q: %w", page, err)
	}

	//At this point parsing succeeded — safe to write headers
	w.WriteHeader(status)

	//Execute the layout
	if err := tmpl.ExecuteTemplate(w, LayoutTemplate, data); err != nil {
		return fmt.Errorf("execute layout %q: %w", LayoutTemplate, err)
	}

	return nil
}

func BadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	if err := Render(w, http.StatusBadRequest, BadRequestTemplate, nil); err != nil {
		logger.L().WithError(err).Error("rendering 400 page failed")
	}
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)

	if err := Render(w, http.StatusNotFound, NotFoundTemplate, nil); err != nil {
		logger.L().WithError(err).Error("rendering 404 page failed")
	}
}

func InternalServerError(w http.ResponseWriter, err error) {
	logger.L().WithError(err).Error("internal server error")

	w.WriteHeader(http.StatusInternalServerError)

	if renderErr := Render(w, http.StatusInternalServerError, InternalServerErrorTemplate, nil); renderErr != nil {
		logger.L().WithError(renderErr).Error("rendering 500 page failed")
	}
}
