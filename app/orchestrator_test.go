package app_test

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/73NN0/foe-hammer/adapters/context"
	hookrunner "github.com/73NN0/foe-hammer/adapters/hook-runner"
	moduleloader "github.com/73NN0/foe-hammer/adapters/module-loader"
	"github.com/73NN0/foe-hammer/app"
	"github.com/73NN0/foe-hammer/domain"
)

func TestOrchestratorComponent(t *testing.T) {
	host := domain.Host{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	rootDir := "../testdata/simple"
	orchestrator := app.NewOrchestrator(
		moduleloader.NewBashLoader(),
		context.NewEnvProvider(),
		hookrunner.NewBashHookRunner(),
		host,
	)

	order, err := orchestrator.GetOrder(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(order) == 0 {
		t.Fatal("Something strange we don't have an order")
	}
}

/*
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/73NN0/foe-hammer/adapters/context"
	"github.com/73NN0/foe-hammer/adapters/hookrunner"
	moduleloader "github.com/73NN0/foe-hammer/adapters/module-loader"
	"github.com/73NN0/foe-hammer/app"
	"github.com/73NN0/foe-hammer/domain"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: foe-hammer <project-dir>")
		os.Exit(1)
	}

	rootDir := os.Args[1]

	host := domain.Host{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	target := domain.Target{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	orchestrator := app.NewOrchestrator(
		moduleloader.NewBashLoader(),
		context.NewEnvProvider(),
		hookrunner.NewBashHookRunner(),
		host,
	)

	order, err := orchestrator.GetOrder(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Build order:")
	for i, name := range order {
		fmt.Printf("  %d. %s\n", i+1, name)
	}
}
*/
