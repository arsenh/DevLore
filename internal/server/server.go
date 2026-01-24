package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arsenh/DevLore/internal/database"
	"github.com/arsenh/DevLore/internal/service"
)

type HTTPServer struct {
	addr      string
	routes    http.Handler
	dbContext *database.DbContext
}

func (s *HTTPServer) notifyContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
}

func NewHTTPServer(addr string) *HTTPServer {
	rateLimiter := NewRateLimiter()
	db := database.NewDbContext()
	db.MustConnect()
	articleService := service.NewArticleService(db)
	userService := service.NewUserService(db)

	return &HTTPServer{
		addr:      addr,
		routes:    NewRoutes(articleService, userService, rateLimiter).GetRoutes(),
		dbContext: db,
	}
}

func (s *HTTPServer) Start() {
	ctx, stop := s.notifyContext()
	defer stop()

	server := http.Server{
		Addr:    s.addr,
		Handler: s.routes,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}

	}()

	<-ctx.Done() // waiting for signals

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.Shutdown(shutdownCtx)
	s.dbContext.Close()
}
