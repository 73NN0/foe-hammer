package context_test

import (
	"path/filepath"
	"testing"

	"github.com/73NN0/foe-hammer/internal/orchestrator/adapters/context"
	"github.com/73NN0/foe-hammer/internal/orchestrator/domain"
)

func TestBuildEnv(t *testing.T) {
	host := domain.Host{OS: "linux", Arch: "amd64"}
	target := domain.Target{OS: "linux", Arch: "amd64"}
	outDir := "/tmp/foe-out"

	mod := &domain.Module{
		Name:    "libb",
		DirPath: "src/libb",
	}

	env := context.NewEnvProvider().BuildEnv(host, target, mod, outDir)
	want := map[string]string{
		"FOE_HOST_OS":     "linux",
		"FOE_HOST_ARCH":   "amd64",
		"FOE_TARGET_OS":   "linux",
		"FOE_TARGET_ARCH": "amd64",
		"FOE_OUTDIR":      outDir,
		"FOE_LIBDIR":      filepath.Join(outDir, "lib"),
		"FOE_BINDIR":      filepath.Join(outDir, "bin"),
		"FOE_OBJDIR":      filepath.Join(outDir, "obj", mod.Name),
		"FOE_SRCDIR":      mod.DirPath,
		"FOE_MODULE_NAME": mod.Name,
	}

	required := []string{
		"FOE_HOST_OS", "FOE_HOST_ARCH", "FOE_TARGET_OS", "FOE_TARGET_ARCH",
		"FOE_OUTDIR", "FOE_LIBDIR", "FOE_BINDIR", "FOE_OBJDIR", "FOE_SRCDIR", "FOE_MODULE_NAME",
	}
	for _, k := range required {
		if _, ok := env[k]; !ok {
			t.Fatalf("missing env key %q", k)
		}
	}

	for k, v := range want {
		got, ok := env[k]
		if !ok {
			t.Fatalf("missing key %q", k)
		}

		if got != v {
			t.Fatalf("key %q: got %q, want %q", k, got, v)
		}
	}

	if len(env) != len(want) {
		t.Fatalf("env has extra/missing keys: got %d keys, want %d keys", len(env), len(want))
	}
}
