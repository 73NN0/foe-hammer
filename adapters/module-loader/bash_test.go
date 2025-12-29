package moduleloader_test

import (
	"testing"

	moduleloader "github.com/73NN0/foe-hammer/adapters/module-loader"
	"github.com/73NN0/foe-hammer/domain"
)

func TestLoadAll(t *testing.T) {
	loader := moduleloader.NewBashLoader()

	modules, err := loader.LoadAll("../../testdata/simple")
	if err != nil {
		t.Fatal(err)
	}

	if len(modules) != 6 {
		t.Errorf("got %d modules, want 6", len(modules))
	}

	var engine *domain.Module
	for _, m := range modules {
		if m.Name == "engine" {
			engine = m
			break
		}
	}

	if engine == nil {
		t.Fatal("engine module not found")
	}

	if len(engine.Depends) != 3 {
		t.Errorf("engine has %d deps, want 3", len(engine.Depends))
	}
}
