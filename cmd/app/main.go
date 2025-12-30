package main

import (
	"github.com/arsenh/DevLore/internal/server"
)

func main() {

	httpServer := server.NewHTTPServer("localhost:8080")
	httpServer.Start()
}
