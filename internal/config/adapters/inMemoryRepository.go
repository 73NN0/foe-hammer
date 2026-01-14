package adapters

import "github.com/73NN0/foe-hammer/internal/config/domain"

type MemoryRepository struct {
	configs map[int]domain.ProjectConfig
}

func NewMemoryRepository(initial []domain.ProjectConfig) *MemoryRepository {
	m := make(map[int]domain.ProjectConfig, len(initial))
	for _, c := range initial {
		m[c.ID] = c
	}
	return &MemoryRepository{configs: m}
}

func (r *MemoryRepository) List() ([]domain.ProjectConfig, error) {
	out := make([]domain.ProjectConfig, 0, len(r.configs))
	for _, v := range r.configs {
		out = append(out, v)
	}
	return out, nil
}

func (r *MemoryRepository) GetByID(id int) (domain.ProjectConfig, error) {
	v, ok := r.configs[id]
	if !ok {
		return domain.ProjectConfig{}, domain.ErrNotFound
	}
	return v, nil
}

var _ domain.Repository = (*MemoryRepository)(nil)
