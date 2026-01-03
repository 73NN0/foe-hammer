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
	tagDeps = "DEPS:"
	tagMake = "MAKE:"
	tagSrcs = "SRCS:"
)

func buildScript(path string) string {
	return `source "` + path + `"
printf '` + tagName + `%s\n' "$pkgname"
printf '` + tagDesc + `%s\n' "$pkgdesc"
printf '` + tagDeps + `%s\n' "${depends[*]}"
printf '` + tagMake + `%s\n' "${makedepends[*]}"
printf '` + tagSrcs + `%s\n' "${source[*]}"`
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

	if err := validateHooks(absPath); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
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
		case strings.HasPrefix(line, tagDeps):
			m.Depends = strings.Fields(strings.TrimPrefix(line, tagDeps))
		case strings.HasPrefix(line, tagMake):
			m.MakeDepends = strings.Fields(strings.TrimPrefix(line, tagMake))
		case strings.HasPrefix(line, tagSrcs):
			m.Sources = strings.Fields(strings.TrimPrefix(line, tagSrcs))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning output: %w", err)
	}
	// validation
	if m.Name == "" {
		return nil, fmt.Errorf("missing pkgname")
	}

	if m.Description == "" {
		return nil, fmt.Errorf("missing pkgdesc")
	}

	if len(m.Sources) == 0 {
		return nil, fmt.Errorf("missing source")
	}

	return m, nil
}

func validateHooks(path string) error {
	script := fmt.Sprintf(`
        source "%s"
        if ! type -t produces &>/dev/null; then
            echo "missing produces()" >&2
            exit 1
        fi
        if ! type -t build &>/dev/null; then
            echo "missing build()" >&2
            exit 1
        fi
    `, path)

	cmd := exec.Command("bash", "-c", script)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("invalid PKGBUILD: %s", strings.TrimSpace(stderr.String()))
	}
	return nil
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

		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if entry.Name() == "bin" || entry.Name() == "build" {
			continue
		}

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
