package domain

import "fmt"

type BuildGraph struct {
	modules  map[string]*Module
	edges    map[string][]string
	indegree map[string]int
}

func NewBuildGraph() *BuildGraph {
	return &BuildGraph{
		modules:  make(map[string]*Module),
		edges:    make(map[string][]string),
		indegree: make(map[string]int),
	}
}

func (g *BuildGraph) AddModule(m *Module) error {
	if _, exists := g.modules[m.Name]; exists {
		return fmt.Errorf("module %s already exists", m.Name)
	}

	g.modules[m.Name] = m
	g.edges[m.Name] = m.Depends

	// init indegree if not done yet

	if _, ok := g.indegree[m.Name]; !ok {
		g.indegree[m.Name] = 0
	}

	for _, dep := range m.Depends {
		g.indegree[dep]++
	}

	return nil
}

func (g *BuildGraph) GetModule(name string) (*Module, error) {
	if module, ok := g.modules[name]; !ok {
		return nil, fmt.Errorf("Module %s doesn't exist", name)
	} else {
		return module, nil
	}
}

// khan algorithm
func (g *BuildGraph) TopoSort() ([]string, error) {
	// TODO !!!!
	return nil, nil
}

func (g *BuildGraph) Validate() error {
	for name, deps := range g.edges {
		for _, dep := range deps {
			if _, exists := g.modules[dep]; !exists {
				return fmt.Errorf("module %s depens on %s, but %s not found", name, dep, dep)
			}
		}
	}
	return nil
}
