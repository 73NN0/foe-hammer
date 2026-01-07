package domain_test

import (
	"testing"

	"github.com/73NN0/foe-hammer/internal/orchestrator/domain"
)

// todo improve this test with cycle limit
func TestTopoSort(t *testing.T) {
	g := domain.NewModuleGraph()

	g.Add(&domain.Module{Name: "libvectors", Depends: []string{}})
	g.Add(&domain.Module{Name: "libcolorbuffer", Depends: []string{}})
	g.Add(&domain.Module{Name: "libmemory", Depends: []string{}})
	g.Add(&domain.Module{Name: "sdl", Depends: []string{}})
	g.Add(&domain.Module{Name: "libkeyboard", Depends: []string{}})
	g.Add(&domain.Module{Name: "libgraphic", Depends: []string{"libvectors", "libcolorbuffer", "libkeyboard"}})
	g.Add(&domain.Module{Name: "engine", Depends: []string{"libgraphic", "libkeyboard", "sdl"}})

	err := g.TopoSort()
	if err != nil {
		t.Fatal(err)
	}
	order := g.Order()
	// verify that every module came after it dependencies
	position := make(map[string]int)
	for i, name := range order {
		position[name] = i
	}

	for _, m := range g.All() {
		for _, dep := range m.Depends {
			if position[dep] > position[m.Name] {
				t.Errorf("%s should come before %s", dep, m.Name)
			}
		}
	}
}
