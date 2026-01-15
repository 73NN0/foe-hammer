package app

import (
	"github.com/73NN0/foe-hammer/internal/config/domain"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

// Create crée une nouvelle config après validation.
// Retourne ErrConfigAlreadyExists si une config existe déjà pour ce path.
func (s *Service) Create(cfg domain.ProjectConfig) error {
	if err := domain.Validate(&cfg); err != nil {
		return err
	}

	// Règle métier : unicité du RootDir
	if _, err := s.repo.GetByPath(cfg.RootDir); err == nil {
		return domain.ErrConfigAlreadyExists
	}

	return s.repo.Create(cfg)
}

// Update met à jour une config existante.
// La config doit exister (vérification par ID).
func (s *Service) Update(cfg domain.ProjectConfig) error {
	if err := domain.Validate(&cfg); err != nil {
		return err
	}

	// Vérifier que la config existe
	existing, err := s.repo.GetByID(cfg.ID)
	if err != nil {
		return domain.ErrConfigNotFound
	}

	// Si le RootDir change, vérifier qu'il n'y a pas de conflit
	if existing.RootDir != cfg.RootDir {
		if _, err := s.repo.GetByPath(cfg.RootDir); err == nil {
			return domain.ErrConfigAlreadyExists
		}
	}

	return s.repo.Update(cfg)
}

// Delete supprime une config par son ID.
func (s *Service) Delete(id int) error {
	// Vérifier que la config existe
	if _, err := s.repo.GetByID(id); err != nil {
		return domain.ErrConfigNotFound
	}
	return s.repo.Delete(id)
}

// GetByID récupère une config par son ID.
func (s *Service) GetByID(id int) (domain.ProjectConfig, error) {
	return s.repo.GetByID(id)
}

// GetByPath récupère une config par son chemin absolu.
func (s *Service) GetByPath(rootDir string) (domain.ProjectConfig, error) {
	return s.repo.GetByPath(rootDir)
}

// List retourne toutes les configs.
func (s *Service) List() ([]domain.ProjectConfig, error) {
	return s.repo.List()
}
