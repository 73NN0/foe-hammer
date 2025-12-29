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
	graph    *domain.BuildGraph
	order    []string
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

func (s *BuildService) Clear() {
	s.graph = nil
	s.order = nil
}

func (s *BuildService) init(rootDir string) ([]string, error) {

	if s.graph != nil {
		// not an error
		return s.order, nil
	}
	// 1. Charger tous les modules
	modules, err := s.loader.LoadAll(rootDir)
	if err != nil {
		return nil, fmt.Errorf("loading modules: %w", err)
	}

	// 2. Construire le graphe
	graph := domain.NewBuildGraph()
	for _, m := range modules {
		if err := graph.AddModule(m); err != nil {
			return nil, fmt.Errorf("adding module %s: %w", m.Name, err)
		}
	}

	// 3. Valider (dépendances manquantes)
	if err := graph.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	// 4. Topo sort
	order, err := graph.TopoSort()
	if err != nil {
		return nil, fmt.Errorf("cycle detected: %w", err)
	}

	s.graph = graph
	return order, nil
}

func (s *BuildService) GetAllExternalDepends(rootDir string) ([]string, error) {
	order, err := s.init(rootDir)
	if err != nil {
		return nil, err
	}

	var allDepends []string
	for _, name := range order {
		m, err := s.graph.GetModule(name)
		if err != nil {
			return nil, err
		}
		allDepends = append(allDepends, m.MakeDepends...)
	}
	return allDepends, nil
}

func (s *BuildService) BuildAll(rootDir string, target domain.Target, outDir string) error {

	order, err := s.init(rootDir)

	if err != nil {
		return err
	}

	for _, name := range order {
		m, err := s.graph.GetModule(name)
		if err != nil {
			return err
		}

		fmt.Printf("Building %s...\n", name)
		if err := s.compiler.Compile(m, target, outDir); err != nil {
			return fmt.Errorf("building %s: %w", name, err)
		}
	}
	return nil
}
