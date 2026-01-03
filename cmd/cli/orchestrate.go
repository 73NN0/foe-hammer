package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/73NN0/foe-hammer/adapters/context"
	hookrunner "github.com/73NN0/foe-hammer/adapters/hook-runner"
	moduleloader "github.com/73NN0/foe-hammer/adapters/module-loader"
	"github.com/73NN0/foe-hammer/adapters/toolchecker"
	"github.com/73NN0/foe-hammer/app"
	"github.com/73NN0/foe-hammer/domain"
)

// this is a basic command
// later this will be split into a build and buildFrom command maybe I don't know

type OrchestrateCommand struct {
	fs *flag.FlagSet
}

func NewOrchestrateCommand() *HelpCommand {
	return &HelpCommand{
		fs: flag.NewFlagSet("help", flag.ExitOnError),
	}
}

func (o *OrchestrateCommand) Name() string           { return "orchestrate" }
func (o *OrchestrateCommand) Description() string    { return "orchestre the build !" }
func (o *OrchestrateCommand) FlagSet() *flag.FlagSet { return o.fs }

func (o *OrchestrateCommand) Run(args []string) error {
	var hostOs, hostArch, targetOs, targetArch, pkgbuildType, rootDir, outDir string

	o.fs.StringVar(&hostOs, "host-os", runtime.GOOS, "host os name")
	o.fs.StringVar(&hostArch, "host-arch", runtime.GOARCH, "host architecture name")
	o.fs.StringVar(&targetOs, "target-os", runtime.GOOS, "target os name")
	o.fs.StringVar(&targetArch, "target-arch", runtime.GOARCH, "host architecture name")
	// note need posix
	o.fs.StringVar(&pkgbuildType, "pkgbuildType", "bash", "type of pkgbuild file (now only bash is available)")
	o.fs.StringVar(&rootDir, "root-dir", ".", "root directory")
	o.fs.StringVar(&outDir, "out-dir", "bin", "output directory")

	if err := o.fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	host := domain.NewHost()
	host.OS = hostOs
	host.Arch = hostArch

	target := domain.NewTarget()
	target.OS = targetOs
	target.Arch = targetArch

	orchestrator := app.NewOrchestrator(
		moduleloader.NewBashLoader(),
		context.NewEnvProvider(),
		hookrunner.NewBashHookRunner(),
		host,
		toolchecker.NewWhichChecker(),
	)

	orchestrator.Load(rootDir)
	orchestrator.SetOutput(outDir)
	orchestrator.Plan(target)
	if err := orchestrator.BuildAll(target); err != nil {
		return err
	}

	return nil
}
