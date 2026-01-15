package domain

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrInvalid             = errors.New("invalid config")
	ErrInvalidRootDir      = errors.New("root_dir is required and must be absolute")
	ErrConfigNotFound      = errors.New("config not found")
	ErrConfigAlreadyExists = errors.New("config already exists for this path")
)

type ProjectConfig struct {
	ID               int      `json:"id"`
	RootDir          string   `json:"root_dir"` // chemin absolu, unique
	ManifestFilename string   `json:"manifest_filename"`
	IgnoreDirs       []string `json:"ignore_dirs"`
	OutDirDefault    string   `json:"out_dir_default"`
}

const (
	defaultManifestName = "PKGBUILD"
	defaultOutDir       = "bin"
)

var defaultIgnoreDirs = []string{"bin", "build", "obj", "node_modules", "vendor"}

// Validate vérifie et complète une config avec les valeurs par défaut.
func Validate(config *ProjectConfig) error {
	rootDir := strings.TrimSpace(config.RootDir)
	if rootDir == "" || !isAbsolutePath(rootDir) {
		return ErrInvalidRootDir
	}
	config.RootDir = rootDir

	if strings.TrimSpace(config.ManifestFilename) == "" {
		config.ManifestFilename = defaultManifestName
	}

	if config.OutDirDefault == "" {
		config.OutDirDefault = defaultOutDir
	}

	if len(config.IgnoreDirs) == 0 {
		config.IgnoreDirs = defaultIgnoreDirs
	}

	return nil
}

func isAbsolutePath(path string) bool {
	p := filepath.Clean(path)

	return filepath.IsAbs(p)
}
