package app

import (
	"fmt"

	"github.com/73NN0/foe-hammer/domain"
	"github.com/73NN0/foe-hammer/ports"
)

type BuildService struct {
	loader   ports.ModuleLoader
	compiler ports.Compiler
	host     domain.Host
}

func NewBuildService(
	loader ports.ModuleLoader,
	compiler ports.Compiler,
	host domain.Host,
) *BuildService {
	return &BuildService{
		loader:   loader,
		compiler: compiler,
		host:     host,
	}
}

func (s *BuildService) BuildAll(rootDir string, target domain.Target, outDir string) error {
	// 1. Charger tous les modules
	modules, err := s.loader.LoadAll(rootDir)
	if err != nil {
		return fmt.Errorf("loading modules: %w", err)
	}

	// 2. Construire le graphe
	graph := domain.NewBuildGraph()
	for _, m := range modules {
		if err := graph.AddModule(m); err != nil {
			return fmt.Errorf("adding module %s: %w", m.Name, err)
		}
	}

	// 3. Valider (dépendances manquantes)
	if err := graph.Validate(); err != nil {
		return fmt.Errorf("validation: %w", err)
	}

	// 4. Topo sort
	order, err := graph.TopoSort()
	if err != nil {
		return fmt.Errorf("cycle detected: %w", err)
	}

	// 5. Compiler dans l'ordre
	for _, name := range order {
		var m *domain.Module

		if module, err := graph.GetModule(name); err != nil {
			return err
		} else {
			(*m) = (*module)
		}
		fmt.Printf("Building %s...\n", name)

		if err := s.compiler.Compile(m, target, outDir); err != nil {
			return fmt.Errorf("building %s: %w", name, err)
		}
	}

	return nil
}
