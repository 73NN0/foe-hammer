package adapters

import (
	"sync"

	"github.com/73NN0/foe-hammer/internal/config/domain"
)

// InMemoryRepository implémente domain.Repository en mémoire.
type InMemoryRepository struct {
	mu      sync.RWMutex
	configs map[int]domain.ProjectConfig
	nextID  int
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		configs: make(map[int]domain.ProjectConfig),
		nextID:  1,
	}
}

func (r *InMemoryRepository) Create(cfg domain.ProjectConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cfg.ID = r.nextID
	r.nextID++
	r.configs[cfg.ID] = cfg
	return nil
}

func (r *InMemoryRepository) Update(cfg domain.ProjectConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.configs[cfg.ID]; !exists {
		return domain.ErrConfigNotFound
	}
	r.configs[cfg.ID] = cfg
	return nil
}

func (r *InMemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.configs[id]; !exists {
		return domain.ErrConfigNotFound
	}
	delete(r.configs, id)
	return nil
}

func (r *InMemoryRepository) GetByID(id int) (domain.ProjectConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cfg, exists := r.configs[id]
	if !exists {
		return domain.ProjectConfig{}, domain.ErrConfigNotFound
	}
	return cfg, nil
}

func (r *InMemoryRepository) GetByPath(rootDir string) (domain.ProjectConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, cfg := range r.configs {
		if cfg.RootDir == rootDir {
			return cfg, nil
		}
	}
	return domain.ProjectConfig{}, domain.ErrConfigNotFound
}

func (r *InMemoryRepository) List() ([]domain.ProjectConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.ProjectConfig, 0, len(r.configs))
	for _, cfg := range r.configs {
		result = append(result, cfg)
	}
	return result, nil
}
