package main

import (
	"os"

	"github.com/73NN0/foe-hammer/internal/config/adapters"
	"github.com/73NN0/foe-hammer/internal/config/app"
	"github.com/73NN0/foe-hammer/internal/config/domain"
	"github.com/73NN0/foe-hammer/internal/config/ports/stdio"
)

func main() {
	// TODO: plus tard: file/env/xdg + merge policy.
	repo := adapters.NewMemoryRepository([]domain.ProjectConfig{
		{
			ID:               1,
			RootDir:          ".",
			ManifestFilename: "PKGBUILD",
			IgnoreDirs:       []string{"bin", "build", "obj", "node_modules", "vendor"},
			OutDirDefault:    "bin",
		},
	})

	svc := app.NewService(repo)
	server := stdio.NewServer(svc)

	if err := server.Serve(os.Stdin, os.Stdout); err != nil {
		// simple pour l’instant
		panic(err)
	}
}
