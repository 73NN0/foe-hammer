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

type OrchestrateCommand struct {
	fs         *flag.FlagSet
	hostOs     string
	hostArch   string
	targetOs   string
	targetArch string
	rootDir    string
	outDir     string
}

func NewOrchestrateCommand() *OrchestrateCommand {
	cmd := &OrchestrateCommand{
		fs: flag.NewFlagSet("orchestrate", flag.ExitOnError),
	}

	// Enregistre les flags ICI, dans le constructeur
	cmd.fs.StringVar(&cmd.hostOs, "host-os", runtime.GOOS, "host os name")
	cmd.fs.StringVar(&cmd.hostArch, "host-arch", runtime.GOARCH, "host architecture name")
	cmd.fs.StringVar(&cmd.targetOs, "target-os", runtime.GOOS, "target os name")
	cmd.fs.StringVar(&cmd.targetArch, "target-arch", runtime.GOARCH, "target architecture name")
	cmd.fs.StringVar(&cmd.rootDir, "root-dir", ".", "root directory")
	cmd.fs.StringVar(&cmd.outDir, "out-dir", "bin", "output directory")

	return cmd
}

func (o *OrchestrateCommand) Name() string           { return "orchestrate" }
func (o *OrchestrateCommand) Description() string    { return "orchestre the build !" }
func (o *OrchestrateCommand) FlagSet() *flag.FlagSet { return o.fs }

func (o *OrchestrateCommand) Run(args []string) error {
	if err := o.fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	host := domain.NewHost()
	host.OS = o.hostOs
	host.Arch = o.hostArch

	target := domain.NewTarget()
	target.OS = o.targetOs
	target.Arch = o.targetArch

	orchestrator := app.NewOrchestrator(
		moduleloader.NewBashLoader(),
		context.NewEnvProvider(),
		hookrunner.NewBashHookRunner(),
		host,
		toolchecker.NewWhichChecker(),
	)

	orchestrator.Load(o.rootDir)
	orchestrator.SetOutput(o.outDir)
	orchestrator.Plan(target)
	if err := orchestrator.BuildAll(target); err != nil {
		return err
	}

	return nil
}
