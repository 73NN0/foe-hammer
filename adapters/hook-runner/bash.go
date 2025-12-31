package hookrunner

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/73NN0/foe-hammer/domain"
)

type BashHookRunner struct{}

func NewBashHookRunner() *BashHookRunner {
	return &BashHookRunner{}
}

func (r *BashHookRunner) Run(module *domain.Module, env map[string]string) error {
	// Source the PKGBUILD and call build()
	script := fmt.Sprintf(`source "%s" && build`, module.Path)

	cmd := exec.Command("bash", "-c", script)
	cmd.Dir = module.DirPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Inject environment variables
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running build hook for %s: %w", module.Name, err)
	}

	return nil
}
