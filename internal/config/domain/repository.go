package domain

// Repository définit les opérations de persistance pour les configs.
type Repository interface {
	Create(cfg ProjectConfig) error
	Update(cfg ProjectConfig) error
	Delete(id int) error
	GetByID(id int) (ProjectConfig, error)
	GetByPath(rootDir string) (ProjectConfig, error)
	List() ([]ProjectConfig, error)
}
