package domain

import (
	"errors"
	"fmt"
)

var (
	ErrGraphCycleDetected     = errors.New("cycle detected")
	ErrGraphModuleNotFound    = errors.New("Module not Found")
	ErrGraphModuleDoesntExist = errors.New("Module doesn't exist")
)

type ModuleGraph struct {
	modules map[string]*Module
	edges   map[string][]string
	order   []string
}

func NewModuleGraph() *ModuleGraph {
	return &ModuleGraph{
		modules: make(map[string]*Module),
		edges:   make(map[string][]string),
	}
}

func (g *ModuleGraph) Add(m *Module) error {
	if _, exists := g.modules[m.Name]; exists {
		return fmt.Errorf("module %s already exists", m.Name)
	}

	g.modules[m.Name] = m
	g.edges[m.Name] = m.Depends

	return nil
}

func (g *ModuleGraph) TopoSort() error {

	_moduleLength := len(g.modules)
	// 1. Construire indegree : combien de deps chaque module attend
	indegree := make(map[string]int, _moduleLength)
	// 2. Construire requiredBy : qui je débloque quand je suis fini
	requiredBy := make(map[string][]string, _moduleLength)

	// Init all known modules indegree to 0

	for name := range g.modules {
		indegree[name] = 0
	}

	// build indegree + required by from the natural representation (module -> dep)

	for modName, m := range g.modules {
		// module -> dep
		indegree[modName] = len(m.Depends)

		for _, dep := range m.Depends {
			// deps-> module ( for khan )
			requiredBy[dep] = append(requiredBy[dep], modName)
		}
	}

	// 3. Queue initiale : modules avec indegree = 0

	queue := make([]string, 0, _moduleLength)
	for name, deg := range indegree {
		if deg == 0 {
			queue = append(queue, name)
		}
	}

	// 4. Boucle : pop, append au résultat, décrémenter les dépendants
	// standard khan ?

	order := make([]string, 0, _moduleLength)
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		order = append(order, n)

		for _, dependent := range requiredBy[n] {
			indegree[dependent]--
			if indegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}

	}
	// 5. Check cycle : si len(result) != len(modules) → erreur
	if len(order) != _moduleLength {
		return ErrGraphCycleDetected
	}

	g.order = order

	return nil
}

func (g *ModuleGraph) Validate() error {
	for name, deps := range g.edges {
		for _, dep := range deps {
			if _, exists := g.modules[dep]; !exists {
				return fmt.Errorf("%w : module %s depends on %s, but %s not found", ErrGraphModuleNotFound, name, dep, dep)
			}
		}
	}
	return nil
}

func (g *ModuleGraph) Get(name string) (*Module, error) {
	if module, ok := g.modules[name]; !ok {
		return nil, fmt.Errorf(" %w, %s", ErrGraphModuleDoesntExist, name)
	} else {
		return module, nil
	}
}

func (g *ModuleGraph) Order() []string {
	return g.order
}

func (g *ModuleGraph) All() (modules []*Module) {
	for _, m := range g.modules {
		modules = append(modules, m)
	}
	return
}

// liba ──→ libb ──→ exe
//
// libX ──→ libY ──→ exe
// Descendants("libb") => ["libb", "exe"]
func (g *ModuleGraph) Descendants(name string) []string {
	// Set des modules à rebuilder
	toRebuild := make(map[string]bool)
	toRebuild[name] = true

	var result []string

	// Parcours dans l'ordre topo
	for _, modName := range g.order {
		// Si ce module est déjà marqué, on l'ajoute au résultat
		if toRebuild[modName] {
			result = append(result, modName)
			continue
		}

		// Sinon, on regarde si une de ses deps est marquée
		m := g.modules[modName]
		for _, dep := range m.Depends {
			if toRebuild[dep] {
				toRebuild[modName] = true
				result = append(result, modName)
				break
			}
		}
	}

	return result
}
