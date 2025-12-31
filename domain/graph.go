package domain

import "fmt"

// src
// https://stackoverflow.com/questions/23531891/how-do-i-succinctly-remove-the-first-element-from-a-slice-in-go?utm_source=chatgpt.com
// but need to search before implementing this kind of queue

// q := make([]string, 0, len(g.modules))
// head := 0

// push := func(x string) { q = append(q, x) }
// pop := func() string {
//     x := q[head]
//     head++
//     // option anti “memory retention” si besoin
//     if head > 1024 && head*2 >= len(q) {
//         q = append([]string(nil), q[head:]...)
//         head = 0
//     }
//     return x
// }
// empty := func() bool { return head >= len(q) }

type BuildGraph struct {
	modules map[string]*Module
	edges   map[string][]string
}

func NewBuildGraph() *BuildGraph {
	return &BuildGraph{
		modules: make(map[string]*Module),
		edges:   make(map[string][]string),
	}
}

func (g *BuildGraph) AddModule(m *Module) error {
	if _, exists := g.modules[m.Name]; exists {
		return fmt.Errorf("module %s already exists", m.Name)
	}

	g.modules[m.Name] = m
	g.edges[m.Name] = m.Depends

	return nil
}

func (g *BuildGraph) GetModule(name string) (*Module, error) {
	if module, ok := g.modules[name]; !ok {
		return nil, fmt.Errorf("Module %s doesn't exist", name)
	} else {
		return module, nil
	}
}

func (g *BuildGraph) TopoSort() ([]string, error) {

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
		return nil, fmt.Errorf("cycle detected (or missing nodes in graph)")
	}

	return order, nil
}

func (g *BuildGraph) Validate() error {
	for name, deps := range g.edges {
		for _, dep := range deps {
			if _, exists := g.modules[dep]; !exists {
				return fmt.Errorf("module %s depends on %s, but %s not found", name, dep, dep)
			}
		}
	}
	return nil
}
