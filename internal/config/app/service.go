package app

import "github.com/73NN0/foe-hammer/internal/config/domain"

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListConfigs() ([]domain.ProjectConfig, error) {
	return s.repo.List()
}

func (s *Service) GetConfig(id int) (domain.ProjectConfig, error) {
	return s.repo.GetByID(id)
}
