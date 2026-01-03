package ports

import "github.com/73NN0/foe-hammer/domain"

// HookRunner executes the build hook of a module
type HookRunner interface {
	Run(module *domain.Module, env map[string]string) error
	Produces(module *domain.Module, env map[string]string) ([]string, error)
}
