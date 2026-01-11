package server

import (
	"net/http"

	"github.com/arsenh/DevLore/internal/logger"
	"github.com/arsenh/DevLore/internal/service"
)

type HTTPServer struct {
	addr   string
	routes http.Handler
}

func NewHTTPServer(addr string) *HTTPServer {
	articleService := service.NewArticleService()
	userService := service.NewUserService()
	return &HTTPServer{
		addr:   addr,
		routes: NewRoutes(articleService, userService).GetRoutes(),
	}
}

func (s *HTTPServer) Start() {
	server := http.Server{
		Addr:    s.addr,
		Handler: s.routes,
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.L().WithError(err).Errorf("unable to start server on addr: %s", s.addr)
	}
}
