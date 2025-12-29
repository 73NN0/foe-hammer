package app_test

import (
	"runtime"
	"testing"

	"github.com/73NN0/foe-hammer/adapters/compiler"
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
	buildService := app.NewBuildService(moduleloader.NewBashLoader(), compiler.NewClangCompiler(), host)

}
