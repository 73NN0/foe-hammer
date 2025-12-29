package moduleloader

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/73NN0/foe-hammer/domain"
)

const (
	tagName = "NAME:"
	tagDesc = "DESC:"
	tagProd = "PROD:"
	tagDeps = "DEPS:"
	tagMake = "MAKE:"
	tagSrcs = "SRCS:"
	tagStgy = "STGY:"
)

func buildScript(path string) string {
	return `source "` + path + `"
printf '` + tagName + `%s\n' "$pkgname"
printf '` + tagDesc + `%s\n' "$pkgdesc"
printf '` + tagProd + `%s\n' "${produces[*]}"
printf '` + tagDeps + `%s\n' "${depends[*]}"
printf '` + tagMake + `%s\n' "${makedepends[*]}"
printf '` + tagSrcs + `%s\n' "${source[*]}"
printf '` + tagStgy + `%s\n' "$strategy"`
}

type BashLoader struct{}

func NewBashLoader() *BashLoader {
	return &BashLoader{}
}

func (l *BashLoader) Load(path string) (*domain.Module, error) {
	absPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("resolving path %s: %w", path, err)
	}
	cmd := exec.Command("bash", "-c", buildScript(absPath))
	cmd.Dir = filepath.Dir(filepath.Clean(path))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("parsing %s: %w\n%s", path, err, stderr.String())
	}
	m, err := l.parseTagged(&stdout)
	if err != nil {
		return nil, err
	}
	m.DirPath = filepath.Dir(absPath)
	m.Path = absPath
	return m, nil
}

func (l *BashLoader) parseTagged(output *bytes.Buffer) (*domain.Module, error) {
	m := &domain.Module{}

	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, tagName):
			m.Name = strings.TrimPrefix(line, tagName)
		case strings.HasPrefix(line, tagDesc):
			m.Description = strings.TrimPrefix(line, tagDesc)
		case strings.HasPrefix(line, tagProd):
			m.Produces = strings.Fields(strings.TrimPrefix(line, tagProd))
		case strings.HasPrefix(line, tagDeps):
			m.Depends = strings.Fields(strings.TrimPrefix(line, tagDeps))
		case strings.HasPrefix(line, tagMake):
			m.MakeDepends = strings.Fields(strings.TrimPrefix(line, tagMake))
		case strings.HasPrefix(line, tagSrcs):
			m.Sources = strings.Fields(strings.TrimPrefix(line, tagSrcs))
		case strings.HasPrefix(line, tagStgy):
			m.Strategy = domain.Strategy(strings.TrimPrefix(line, tagStgy))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning output: %w", err)
	}

	if m.Name == "" {
		return nil, fmt.Errorf("missing pkgname")
	}

	return m, nil
}

func (l *BashLoader) LoadAll(rootDir string) ([]*domain.Module, error) {
	var modules []*domain.Module

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// improve this :
		// it's hardcoded
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if entry.Name() == "bin" || entry.Name() == "build" {
			continue
		}
		// Vérifier si PKGBUILD existe (insensible à la casse)
		path, found := l.findManifest(filepath.Join(rootDir, entry.Name()))
		if !found {
			continue
		}

		m, err := l.Load(path)
		if err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}

	return modules, nil
}

func (l *BashLoader) findManifest(dir string) (string, bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", false
	}

	for _, entry := range entries {
		if strings.EqualFold(entry.Name(), "PKGBUILD") {
			return filepath.Join(dir, entry.Name()), true
		}
	}

	return "", false
}
