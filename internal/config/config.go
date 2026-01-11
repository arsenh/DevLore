package config

import "path/filepath"

var (
	// TemplatesDir HTML templates directory
	TemplatesDir     = filepath.Join("internal", "templates")
	StaticDir        = "static"
	PublicDir        = "public"
	MainHTMLTemplate = filepath.Join(TemplatesDir, "layout.html")
	JWTSecretKey     = []byte("super-secret")
)
