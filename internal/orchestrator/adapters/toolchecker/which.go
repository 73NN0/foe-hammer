package toolchecker

import (
	"fmt"
	"os/exec"

	"github.com/73NN0/foe-hammer/internal/orchestrator/domain"
)

type WhichChecker struct{}

func NewWhichChecker() *WhichChecker {
	return &WhichChecker{}
}

func (c *WhichChecker) Check(tool string) error {
	_, err := exec.LookPath(tool)
	if err != nil {
		return fmt.Errorf("%s not found in PATH", tool)
	}
	return nil
}

func (c *WhichChecker) Suggest(tool string, host domain.Host) string {
	// very basic
	return fmt.Sprintf("install %s", tool)
}
