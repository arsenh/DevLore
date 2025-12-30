package server

import (
	"fmt"
	"net/http"
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
		Handler: s.Routes,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("unable to start server on addr:", s.Addr)
	}
}
