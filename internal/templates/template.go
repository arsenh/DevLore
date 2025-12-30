package templates

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/arsenh/DevLore/internal/config"
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
	// TODO: add correct logging
	fmt.Println("templates are parsed successfully.")
}

func RenderHTML(w http.ResponseWriter, page string, data interface{}) {
	// Clone the base template.
	tmpl, err := base.Clone()
	if err != nil {
		http.Error(w, "failed to clone template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the specific page template into the clone.
	fullPagePath := filepath.Join(config.TemplatesDir, page)
	fmt.Println("fullPagePath:", fullPagePath)
	_, err = tmpl.ParseFiles(fullPagePath)
	if err != nil {
		http.Error(w, "failed to parse page template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the layout with the data.
	err = tmpl.ExecuteTemplate(w, Layout, data)
	if err != nil {
		http.Error(w, "failed to execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
