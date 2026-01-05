package templates

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/arsenh/DevLore/internal/config"
	"github.com/arsenh/DevLore/internal/logger"
)

const (
	LayoutTemplate              = "layout"
	DashboardTemplate           = "dashboard.html"
	NewArticleTemplate          = "new-article.html"
	EditArticleTemplate         = "edit-article.html"
	ViewArticleTemplate         = "view-article.html"
	Search                      = "search.html"
	NotFoundTemplate            = "404.html"
	BadRequestTemplate          = "400.html"
	InternalServerErrorTemplate = "500.html"
)

var base *template.Template

func init() {
	basePath := filepath.Join(config.MainHTMLTemplate)
	base = template.Must(template.ParseFiles(basePath))
	logger.L().WithField("path", basePath).Info("rendering template")
	logger.L().Info("templates are parsed successfully")
}

func Render(w http.ResponseWriter, status int, page string, data interface{}) error {
	//Clone the base layout template
	tmpl, err := base.Clone()
	if err != nil {
		return fmt.Errorf("clone base template %q: %w", base.Name(), err)
	}

	//Parse the page template into the clone
	fullPagePath := filepath.Join(config.TemplatesDir, page)

	if _, err := tmpl.ParseFiles(fullPagePath); err != nil {
		return fmt.Errorf("parse page template %q: %w", fullPagePath, err)
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
