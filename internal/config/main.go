package main

import (
	"os"

	"github.com/73NN0/foe-hammer/internal/common/stdio"
	"github.com/73NN0/foe-hammer/internal/config/adapters"
	"github.com/73NN0/foe-hammer/internal/config/app"
	"github.com/73NN0/foe-hammer/internal/config/ports"
)

func main() {
	// TODO: plus tard: file/env/xdg + merge policy.

	repo := adapters.NewInMemoryRepository()

	service := app.NewService(repo)
	server := stdio.NewServer(stdio.ServerConfig{
		Topic:   "config.command",
		Handler: ports.NewStdioConfigHandler(service),
	})

	if err := server.Serve(os.Stdin, os.Stdout); err != nil {
		// simple pour l’instant
		panic(err)
	}
}
