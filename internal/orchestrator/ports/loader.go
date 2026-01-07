package ports

import "github.com/73NN0/foe-hammer/internal/orchestrator/domain"

type ModuleLoader interface {
	LoadAll(rootDir string) ([]*domain.Module, error)
	Load(path string) (*domain.Module, error)
}
