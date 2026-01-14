package domain

import "errors"

var ErrNotFound = errors.New("config not found")

type Repository interface {
	List() ([]ProjectConfig, error)
	GetByID(int) (ProjectConfig, error)
}
