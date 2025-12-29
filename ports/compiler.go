package ports

import "github.com/73NN0/foe-hammer/domain"

// compiler knows how transform sources files to artefacts

type Compiler interface {
	Compile(m *domain.Module, target domain.Target, outdir string) error
	// return true if the compiler knows how to handle this target
	CanHandle(host domain.Host, target domain.Target) bool
}
