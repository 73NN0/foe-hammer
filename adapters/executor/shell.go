package executor

import (
	"fmt"
	"io"
	"os/exec"
)

type Shell struct{}

func NewShell() *Shell {
	return &Shell{}
}

func (e *Shell) Run(cmd string, args []string, workDir string, stdout, stderr io.Writer) error {
	c := exec.Command(cmd, args...)
	c.Dir = workDir
	c.Stdout = stdout
	c.Stderr = stderr

	if err := c.Run(); err != nil {
		return fmt.Errorf("%s: %w", cmd, err)
	}
	return nil
}
