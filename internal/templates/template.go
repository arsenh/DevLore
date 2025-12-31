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
	Layout      = "layout"
	Dashboard   = "dashboard.html"
	NewArticle  = "new-article.html"
	EditArticle = "edit-article.html"
	ViewArticle = "view-article.html"
	Search      = "search.html"
)

var base *template.Template

func init() {
	basePath := filepath.Join(config.MainHTMLTemplate)
	base = template.Must(template.ParseFiles(basePath))
	logger.L().WithField("path", basePath).Info("rendering template")
	logger.L().Info("templates are parsed successfully")
}

func RenderHTML(w http.ResponseWriter, page string, data interface{}) {
	// Clone the base template.
	tmpl, err := base.Clone()
	if err != nil {
		msg := fmt.Sprintf("failed to clone base template: %s", base.Name())
		logger.L().WithError(err).Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Parse the specific page template into the clone.
	fullPagePath := filepath.Join(config.TemplatesDir, page)
	_, err = tmpl.ParseFiles(fullPagePath)
	if err != nil {
		msg := fmt.Sprintf("failed to parse page template: %s", fullPagePath)
		logger.L().WithError(err).Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Execute the layout with the data.
	err = tmpl.ExecuteTemplate(w, Layout, data)
	if err != nil {
		msg := fmt.Sprintf("failed to execute template: %s", Layout)
		logger.L().WithError(err).Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}
