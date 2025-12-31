package server

import (
	"net/http"

	"github.com/arsenh/DevLore/internal/logger"
	"github.com/arsenh/DevLore/internal/middleware"
)

type HTTPServer struct {
	Addr   string
	Routes *http.ServeMux
}

func NewHTTPServer(addr string) *HTTPServer {
	return &HTTPServer{
		Addr:   addr,
		Routes: GetRoutes(),
	}
}

func (s *HTTPServer) Start() {
	server := http.Server{
		Addr:    s.Addr,
		Handler: middleware.Logging(s.Routes),
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.L().WithError(err).Errorf("unable to start server on addr: %s", s.Addr)
	}
}
