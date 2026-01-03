package app_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/73NN0/foe-hammer/adapters/context"
	hookrunner "github.com/73NN0/foe-hammer/adapters/hook-runner"
	moduleloader "github.com/73NN0/foe-hammer/adapters/module-loader"
	"github.com/73NN0/foe-hammer/adapters/toolchecker"
	"github.com/73NN0/foe-hammer/app"
	"github.com/73NN0/foe-hammer/domain"
)

const (
	simplePath string = "../testdata/simple"
	cyclePath  string = "../testdata/cycle"
	plan9Path  string = "../testdata/plan9"
)

func TestOrchestrator(t *testing.T) {
	host := domain.Host{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
	target := domain.Target{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	tests := []struct {
		name         string
		rootDir      string
		wantLoadErr  bool
		wantPlanErr  bool
		wantBuildErr bool
		checkOrder   func(t *testing.T, order []string)
		checkBuild   func(t *testing.T, outDir string, modules []*domain.Module)
	}{
		{
			name:    "happy path - simple chain",
			rootDir: simplePath,
			checkOrder: func(t *testing.T, order []string) {
				// app/engine doit être en dernier
				last := order[len(order)-1]
				if last != "app" && last != "engine" {
					t.Errorf("expected app or engine last, got %s", last)
				}
			},
			checkBuild: func(t *testing.T, outDir string, modules []*domain.Module) {
				for _, m := range modules {
					for _, produce := range m.Produces {
						path := filepath.Join(outDir, produce)
						if _, err := os.Stat(path); os.IsNotExist(err) {
							t.Errorf("%s: expected %s to exist", m.Name, path)
						}
					}
				}
			},
		},
		{
			name:        "cycle detection",
			rootDir:     cyclePath,
			wantLoadErr: true,
		},
		{
			name:    "plan9 - platform abstraction",
			rootDir: plan9Path,
			checkOrder: func(t *testing.T, order []string) {
				// app should be last
				last := order[len(order)-1]
				if last != "app" {
					t.Errorf("expected app last, got %s", last)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outDir := t.TempDir()

			orchestrator := app.NewOrchestrator(
				moduleloader.NewBashLoader(),
				context.NewEnvProvider(),
				hookrunner.NewBashHookRunner(),
				host,
				toolchecker.NewWhichChecker(),
			)

			// Load
			err := orchestrator.Load(tt.rootDir)
			if tt.wantLoadErr {
				if err == nil {
					t.Fatal("expected Load error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Load: %v", err)
			}

			// Check order
			if tt.checkOrder != nil {
				tt.checkOrder(t, orchestrator.Order())
			}

			// Plan
			orchestrator.SetOutput(outDir)
			t.Logf("outDir = %s", outDir)
			err = orchestrator.Plan(target)
			if tt.wantPlanErr {
				if err == nil {
					t.Fatal("expected Plan error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Plan: %v", err)
			}

			// Build (optionnel selon le test)
			if tt.checkBuild != nil {
				err = orchestrator.BuildAll(target)
				if tt.wantBuildErr {
					if err == nil {
						t.Fatal("expected BuildAll error, got nil")
					}
					return
				}
				if err != nil {
					t.Fatalf("BuildAll: %v", err)
				}

				tt.checkBuild(t, outDir, orchestrator.All())
			}
		})
	}
}

func TestPlan(t *testing.T) {

}
