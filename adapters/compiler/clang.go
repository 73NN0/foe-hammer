package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/73NN0/foe-hammer/domain"
	"github.com/73NN0/foe-hammer/ports"
)

type ClangCompiler struct {
	executor ports.Executor
	host     domain.Host
}

func NewClangCompiler(executor ports.Executor, host domain.Host) *ClangCompiler {
	return &ClangCompiler{executor: executor, host: host}
}

func (c *ClangCompiler) compileToObjects(m *domain.Module, outDir string, config domain.CrossConfig) ([]string, error) {
	var objects []string
	for _, src := range m.Sources {
		obj := strings.TrimSuffix(filepath.Base(src), ".c") + ".o"
		objPath := filepath.Join(outDir, "obj", m.Name, obj)

		args := []string{"-c", src, "-o", objPath}
		args = append(args, config.Flags...)

		if err := c.executor.Run(config.Compiler, args, m.DirPath, os.Stdout, os.Stderr); err != nil {
			return nil, err
		}

		objects = append(objects, objPath)
	}

	return objects, nil
}

func (c *ClangCompiler) Compile(m *domain.Module, target domain.Target, outDir string) error {
	config := domain.GetCrossConfig(c.host, target)

	// 1. Compiler les sources en .o
	objects, err := c.compileToObjects(m, outDir, config)
	if err != nil {
		return fmt.Errorf("compiling: %w", err)
	}

	// 2. Produire l'artefact déclaré
	for _, prod := range m.Produces {
		if strings.HasSuffix(prod, ".a") {
			// Archive
			libPath := filepath.Join(outDir, prod)
			args := append([]string{"rcs", libPath}, objects...)
			if err := c.executor.Run("ar", args, ".", os.Stdout, os.Stderr); err != nil {
				return err
			}
			// TODO : Improve this hardcoded directory
		} else if strings.Contains(prod, "bin/") {
			// Link executable
			exePath := filepath.Join(outDir, prod)
			args := append([]string{"-o", exePath}, objects...)
			args = append(args, config.Flags...)
			// TODO: ajouter les libs des dépendances
			if err := c.executor.Run(config.Compiler, args, ".", os.Stdout, os.Stderr); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ClangCompiler) CanHandle(host domain.Host, target domain.Target) bool {
	_, err := exec.LookPath("clang")
	return err == nil
}
