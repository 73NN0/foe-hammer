package ports

import "github.com/73NN0/foe-hammer/internal/orchestrator/domain"

// ContextProvider builds the environment variables for hook execution
type ContextProvider interface {
	BuildEnv(host domain.Host, target domain.Target, module *domain.Module, outDir string) map[string]string
}
