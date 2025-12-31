package app

import (
	"fmt"

	"github.com/73NN0/foe-hammer/domain"
	"github.com/73NN0/foe-hammer/ports"
)

type Orchestrator struct {
	loader  ports.ModuleLoader
	context ports.ContextProvider
	runner  ports.HookRunner
	host    domain.Host
	graph   *domain.BuildGraph
	order   []string
}

func NewOrchestrator(
	loader ports.ModuleLoader,
	context ports.ContextProvider,
	runner ports.HookRunner,
	host domain.Host,
) *Orchestrator {
	return &Orchestrator{
		loader:  loader,
		context: context,
		runner:  runner,
		host:    host,
	}
}

func (o *Orchestrator) Clear() {
	o.graph = nil
	o.order = nil
}

func (o *Orchestrator) init(rootDir string) ([]string, error) {
	if o.graph != nil {
		return o.order, nil
	}

	// 1. Load all modules
	modules, err := o.loader.LoadAll(rootDir)
	if err != nil {
		return nil, fmt.Errorf("loading modules: %w", err)
	}

	// 2. Build the graph
	graph := domain.NewBuildGraph()
	for _, m := range modules {
		if err := graph.AddModule(m); err != nil {
			return nil, fmt.Errorf("adding module %s: %w", m.Name, err)
		}
	}

	// 3. Validate (missing dependencies)
	if err := graph.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	// 4. Topo sort
	order, err := graph.TopoSort()
	if err != nil {
		return nil, fmt.Errorf("cycle detected: %w", err)
	}

	o.graph = graph
	o.order = order
	return order, nil
}

func (o *Orchestrator) GetOrder(rootDir string) ([]string, error) {
	return o.init(rootDir)
}

func (o *Orchestrator) GetAllExternalDepends(rootDir string) ([]string, error) {
	order, err := o.init(rootDir)
	if err != nil {
		return nil, err
	}

	var allDepends []string
	for _, name := range order {
		m, err := o.graph.GetModule(name)
		if err != nil {
			return nil, err
		}
		allDepends = append(allDepends, m.MakeDepends...)
	}
	return allDepends, nil
}

func (o *Orchestrator) BuildAll(rootDir string, target domain.Target, outDir string) error {
	order, err := o.init(rootDir)
	if err != nil {
		return err
	}

	for _, name := range order {
		m, err := o.graph.GetModule(name)
		if err != nil {
			return err
		}

		fmt.Printf("Building %s...\n", name)

		// Build environment for this module
		env := o.context.BuildEnv(o.host, target, m, outDir)

		// Run the build hook
		if err := o.runner.Run(m, env); err != nil {
			return fmt.Errorf("building %s: %w", name, err)
		}
	}

	return nil
}
