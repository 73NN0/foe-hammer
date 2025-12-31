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

func objectPath(absOutDir, moduleName, srcFile string) string {
	obj := strings.TrimSuffix(filepath.Base(srcFile), ".c") + ".o"
	return filepath.Join(absOutDir, "obj", moduleName, obj)
}

func compileArgs(src, objPath string, flags []string) []string {
	args := []string{"-c", src, "-o", objPath}
	return append(args, flags...)
}

// Fonctions I/O (petites, une seule responsabilité)
func resolveAbsPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolving path %s: %w", path, err)
	}
	return abs, nil
}

func ensureDir(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating dir for %s: %w", path, err)
	}
	return nil
}

func (c *ClangCompiler) runCompile(compiler string, args []string, workDir string) error {
	if err := c.executor.Run(compiler, args, workDir, os.Stdout, os.Stderr); err != nil {
		return fmt.Errorf("%s: %w", compiler, err)
	}
	return nil
}

func (c *ClangCompiler) compileToObjects(m *domain.Module, outDir string, config domain.CrossConfig) ([]string, error) {
	absOutDir, err := resolveAbsPath(outDir)
	if err != nil {
		return nil, err
	}

	objects := make([]string, 0, len(m.Sources))

	for _, src := range m.Sources {
		objPath, err := c.compileOneObject(src, absOutDir, m.Name, m.DirPath, config)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objPath)
	}

	return objects, nil
}

func (c *ClangCompiler) compileOneObject(src, absOutDir, moduleName, workDir string, config domain.CrossConfig) (string, error) {
	objPath := objectPath(absOutDir, moduleName, src)

	if err := ensureDir(objPath); err != nil {
		return "", err
	}

	args := compileArgs(src, objPath, config.Flags)

	if err := c.runCompile(config.Compiler, args, workDir); err != nil {
		return "", err
	}

	return objPath, nil
}

// -----

// Fonctions pures
func archiveArgs(libPath string, objects []string) []string {
	return append([]string{"rcs", libPath}, objects...)
}

func linkArgs(exePath string, objects []string, flags []string) []string {
	args := append([]string{"-o", exePath}, objects...)
	return append(args, flags...)
}

func artifactPath(outDir, prod string) string {
	return filepath.Join(outDir, prod)
}

func isLibrary(prod string) bool {
	return strings.HasSuffix(prod, ".a")
}

func isExecutable(prod string) bool {
	return strings.Contains(prod, "bin/")
}

// Fonctions I/O
func (c *ClangCompiler) runArchive(libPath string, objects []string) error {
	if err := ensureDir(libPath); err != nil {
		return err
	}
	args := archiveArgs(libPath, objects)
	return c.executor.Run("ar", args, ".", os.Stdout, os.Stderr)
}

func (c *ClangCompiler) runLink(exePath string, objects []string, config domain.CrossConfig) error {
	if err := ensureDir(exePath); err != nil {
		return err
	}
	args := linkArgs(exePath, objects, config.Flags)
	return c.executor.Run(config.Compiler, args, ".", os.Stdout, os.Stderr)
}

func (c *ClangCompiler) produceArtifact(prod, outDir string, objects []string, config domain.CrossConfig) error {
	path := artifactPath(outDir, prod)

	switch {
	case isLibrary(prod):
		return c.runArchive(path, objects)
	case isExecutable(prod):
		return c.runLink(path, objects, config)
	default:
		return fmt.Errorf("unknown artifact type: %s", prod)
	}
}

// Orchestrateur
func (c *ClangCompiler) Compile(m *domain.Module, target domain.Target, outDir string) error {
	config := domain.GetCrossConfig(c.host, target)

	objects, err := c.compileToObjects(m, outDir, config)
	if err != nil {
		return fmt.Errorf("compiling: %w", err)
	}

	for _, prod := range m.Produces {
		if err := c.produceArtifact(prod, outDir, objects, config); err != nil {
			return fmt.Errorf("producing %s: %w", prod, err)
		}
	}

	return nil
}

func (c *ClangCompiler) CanHandle(host domain.Host, target domain.Target) bool {
	_, err := exec.LookPath("clang")
	return err == nil
}
