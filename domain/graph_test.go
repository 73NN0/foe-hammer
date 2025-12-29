package domain

import "testing"

// todo improve this test with cycle limit
func TestTopoSort(t *testing.T) {
	g := NewBuildGraph()

	g.AddModule(&Module{Name: "libvectors", Depends: []string{}})
	g.AddModule(&Module{Name: "libcolorbuffer", Depends: []string{}})
	g.AddModule(&Module{Name: "libmemory", Depends: []string{}})
	g.AddModule(&Module{Name: "sdl", Depends: []string{}})
	g.AddModule(&Module{Name: "libkeyboard", Depends: []string{}})
	g.AddModule(&Module{Name: "libgraphic", Depends: []string{"libvectors", "libcolorbuffer", "libkeyboard"}})
	g.AddModule(&Module{Name: "engine", Depends: []string{"libgraphic", "libkeyboard", "sdl"}})

	order, err := g.TopoSort()
	if err != nil {
		t.Fatal(err)
	}

	// verify that every module came after it dependencies
	position := make(map[string]int)
	for i, name := range order {
		position[name] = i
	}

	for _, m := range g.modules {
		for _, dep := range m.Depends {
			if position[dep] > position[m.Name] {
				t.Errorf("%s should come before %s", dep, m.Name)
			}
		}
	}
}
