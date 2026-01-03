package app

import (
	"fmt"

	"github.com/73NN0/foe-hammer/domain"
	"github.com/73NN0/foe-hammer/ports"
)

// Orchestrator coordinates the build of modules.
// It loads modules, resolves dependencies, and executes build hooks.
//
// Typical usage:
//
//	o := app.NewOrchestrator(loader, context, runner, host, checker)
//	o.Load("./project")
//	o.SetOutput("./build")
//	o.Plan(target)
//	o.BuildAll(target)
type Orchestrator struct {
	loader  ports.ModuleLoader
	context ports.ContextProvider
	runner  ports.HookRunner
	checker ports.ToolChecker // CanBuild
	host    domain.Host
	graph   *domain.ModuleGraph
	rootDir string
	outDir  string
}

func NewOrchestrator(
	loader ports.ModuleLoader,
	context ports.ContextProvider,
	runner ports.HookRunner,
	host domain.Host,
	checker ports.ToolChecker,
) *Orchestrator {
	return &Orchestrator{
		loader:  loader,
		context: context,
		runner:  runner,
		host:    host,
		checker: checker,
	}
}

// Load scans rootDir for modules, builds the dependency graph,
// validates dependencies, and computes the topological order.
func (o *Orchestrator) Load(rootDir string) error {
	// 1. Load all modules
	modules, err := o.loader.LoadAll(rootDir)
	if err != nil {
		return fmt.Errorf("loading modules: %w", err)
	}

	// 2. Build the graph
	graph := domain.NewModuleGraph()
	for _, m := range modules {
		if err := graph.Add(m); err != nil {
			return fmt.Errorf("adding module %s: %w", m.Name, err)
		}
	}

	// 3. Validate (missing dependencies)
	if err := graph.Validate(); err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	// 4. Topo sort
	if err := graph.TopoSort(); err != nil {
		return err
	}

	o.graph = graph

	return nil
}

// SetOutput sets the output directory for build artifacts.
// Must be called before Plan.
func (o *Orchestrator) SetOutput(outDir string) {
	o.outDir = outDir
}

// All returns all loaded modules directly from the internal graph so there is no topoligical order.
// Requires: Load must be called first.
func (o *Orchestrator) All() []*domain.Module {
	return o.graph.All()
}

// Order returns the topological build order.
// Order result is not deterministic (flemme)
// Requires: Load must be called first.
func (o *Orchestrator) Order() []string {
	return o.graph.Order()
}

// CanBuild checks if all external tools (makedepends) are available for a module.
// Module dependencies check is assured by the graph construction when loading the orchestrator
func (o *Orchestrator) CanBuild(name string) error {
	m, err := o.graph.Get(name)
	if err != nil {
		return err
	}

	for _, tool := range m.MakeDepends {
		if err := o.checker.Check(tool); err != nil {
			return fmt.Errorf("missing tool %s for %s: %w", tool, name, err)
		}
	}

	return nil
}

// Build builds a single module.
// Requires: Plan must be called first.
func (o *Orchestrator) Build(name string, target domain.Target) error {
	fmt.Printf("DEBUG: outDir = %q\n", o.outDir)
	m, err := o.graph.Get(name)
	if err != nil {
		return err
	}

	// verify tools
	if err := o.CanBuild(name); err != nil {
		return err
	}

	// prepare env
	env := o.context.BuildEnv(o.host, target, m, o.outDir)

	// execute hook
	if err := o.runner.Run(m, env); err != nil {
		return fmt.Errorf("building %s: %w", name, err)
	}

	return nil
}

// BuildFrom builds a module and all its descendants (modules that depend on it).
// Requires: Plan must be called first.
func (o *Orchestrator) BuildFrom(name string, target domain.Target) error {
	// Descendants retourne [name, ...ceux qui dépendent de name] dans l'ordre topo
	toBuild := o.graph.Descendants(name)

	for _, modName := range toBuild {
		fmt.Printf("Building %s...\n", modName)
		if err := o.Build(modName, target); err != nil {
			return err
		}
	}

	return nil
}

// BuildAll builds all modules in topological order.
// Requires: Plan must be called first.
func (o *Orchestrator) BuildAll(target domain.Target) error {
	order := o.graph.Order()

	for _, name := range order {
		fmt.Printf("Building %s...\n", name)
		if err := o.Build(name, target); err != nil {
			return err
		}
	}

	return nil
}

// Plan resolves what each module will produce for the given target.
// Must be called after Load and SetOutput, and before any Build method.
func (o *Orchestrator) Plan(target domain.Target) error {
	for _, m := range o.graph.All() {
		env := o.context.BuildEnv(o.host, target, m, o.outDir)

		produces, err := o.runner.Produces(m, env)
		if err != nil {
			return fmt.Errorf("resolving produces for %s: %w", m.Name, err)
		}

		m.Produces = produces
	}
	return nil
}
