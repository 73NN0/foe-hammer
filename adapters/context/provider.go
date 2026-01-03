package context

import (
	"path/filepath"

	"github.com/73NN0/foe-hammer/domain"
)

type EnvProvider struct{}

func NewEnvProvider() *EnvProvider {
	return &EnvProvider{}
}

// module can't be null or hasn't have a name the module loader must garanty it
func (p *EnvProvider) BuildEnv(host domain.Host, target domain.Target, module *domain.Module, outDir string) map[string]string {

	return map[string]string{
		"FOE_HOST_OS":     host.OS,
		"FOE_HOST_ARCH":   host.Arch,
		"FOE_TARGET_OS":   target.OS,
		"FOE_TARGET_ARCH": target.Arch,
		"FOE_OUTDIR":      outDir,
		"FOE_LIBDIR":      filepath.Join(outDir, "lib"),
		"FOE_BINDIR":      filepath.Join(outDir, "bin"),
		"FOE_OBJDIR":      filepath.Join(outDir, "obj", module.Name),
		"FOE_SRCDIR":      module.DirPath,
		"FOE_MODULE_NAME": module.Name,
	}
}
