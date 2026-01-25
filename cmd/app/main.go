package main

import (
	"github.com/arsenh/DevLore/internal/server"
)

func main() {
	httpServer := server.NewHTTPServer()
	httpServer.Start()
}
