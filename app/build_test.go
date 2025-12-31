package app_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/73NN0/foe-hammer/adapters/compiler"
	"github.com/73NN0/foe-hammer/adapters/executor"
	moduleloader "github.com/73NN0/foe-hammer/adapters/module-loader"
	"github.com/73NN0/foe-hammer/app"
	"github.com/73NN0/foe-hammer/domain"
)

func TestBuild(t *testing.T) {
	// setup
	host := domain.Host{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	target := domain.Target{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	binDirectory := "../testdata/simple/bin"

	buildService := app.NewBuildService(moduleloader.NewBashLoader(), compiler.NewClangCompiler(executor.NewShell(), host))

	if err := buildService.BuildAll("../testdata/simple", target, binDirectory); err != nil {
		os.RemoveAll(binDirectory)
		t.Fatal(err)
	}

	os.RemoveAll(binDirectory)
}
