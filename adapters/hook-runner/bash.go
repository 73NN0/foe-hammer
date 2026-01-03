package hookrunner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/73NN0/foe-hammer/domain"
)

type BashHookRunner struct{}

func NewBashHookRunner() *BashHookRunner {
	return &BashHookRunner{}
}

func (r *BashHookRunner) Run(module *domain.Module, env map[string]string) error {
	script := fmt.Sprintf(`source "%s" && build`, module.Path)

	cmd := createCmd(script, module.DirPath, os.Stdout)
	injectEnvv(cmd, env)
	return execute(cmd)
}

func (r *BashHookRunner) Produces(module *domain.Module, env map[string]string) ([]string, error) {
	script := fmt.Sprintf(`source "%s" && produces`, module.Path)

	var stdout bytes.Buffer
	cmd := createCmd(script, module.DirPath, &stdout)
	injectEnvv(cmd, env)

	if err := execute(cmd); err != nil {
		return nil, err
	}

	// Parse output : une ligne = un produce
	var produces []string
	for _, line := range strings.Split(strings.TrimSpace(stdout.String()), "\n") {
		if line != "" {
			produces = append(produces, line)
		}
	}
	return produces, nil
}

func createCmd(script, DirPath string, stdout io.Writer) *exec.Cmd {
	cmd := exec.Command("bash", "-c", script)
	cmd.Dir = DirPath
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func injectEnvv(cmd *exec.Cmd, env map[string]string) {

	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
}

func execute(cmd *exec.Cmd) error {
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("executing %s: %w", cmd.Args, err)
	}
	return nil
}
