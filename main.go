// main.go

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/73NN0/foe-hammer/adapters/clang"
	"github.com/73NN0/foe-hammer/adapters/pkgbuild"
	"github.com/73NN0/foe-hammer/adapters/shell"
	"github.com/73NN0/foe-hammer/app"
	"github.com/73NN0/foe-hammer/domain"
)

func main() {
	// Détecter le host
	host := domain.Host{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	// Target (pour l'instant = host, plus tard via flags)
	target := domain.Target{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	// Câbler les adapters
	executor := shell.NewExecutor()
	loader := pkgbuild.NewLoader()
	compiler := clang.NewCompiler(executor, host)

	// Créer le service
	buildService := app.NewBuildService(loader, compiler, host)

	// Exécuter
	if err := buildService.BuildAll(".", target, "bin"); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
