package moduleloader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid PKGBUILD",
			path:    "../../testdata/loader/valid/PKGBUILD",
			wantErr: false,
		},
		{
			name:        "missing pkgname",
			path:        "../../testdata/loader/missing-pkgname/PKGBUILD",
			wantErr:     true,
			errContains: "missing pkgname",
		},
		{
			name:        "missing pkgdesc",
			path:        "../../testdata/loader/missing-pkgdesc/PKGBUILD",
			wantErr:     true,
			errContains: "missing pkgdesc",
		},
		{
			name:        "missing source",
			path:        "../../testdata/loader/missing-source/PKGBUILD",
			wantErr:     true,
			errContains: "missing source",
		},
		{
			name:        "missing produces hook",
			path:        "../../testdata/loader/missing-produces-hook/PKGBUILD",
			wantErr:     true,
			errContains: "missing produces()",
		},
		{
			name:        "missing build hook",
			path:        "../../testdata/loader/missing-build-hook/PKGBUILD",
			wantErr:     true,
			errContains: "missing build()",
		},
	}

	loader := NewBashLoader()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := loader.Load(tt.path)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if m == nil {
				t.Fatal("expected module, got nil")
			}

			if m.Name == "" {
				t.Error("module name is empty")
			}
		})
	}
}

func TestValidateHooks(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid hooks",
			path:    "../../testdata/loader/valid/PKGBUILD",
			wantErr: false,
		},
		{
			name:        "missing produces",
			path:        "../../testdata/loader/missing-produces-hook/PKGBUILD",
			wantErr:     true,
			errContains: "missing produces()",
		},
		{
			name:        "missing build",
			path:        "../../testdata/loader/missing-build-hook/PKGBUILD",
			wantErr:     true,
			errContains: "missing build()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHooks(tt.path)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestLoadAllRecursive(t *testing.T) {
	loader := NewBashLoader()

	modules, err := loader.LoadAll("../../testdata/plan9")
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	expectedModules := map[string]bool{
		"platform-stub": false,
		"platform-unix": false,
		"libcore":       false,
		"app":           false,
	}

	for _, m := range modules {
		if _, ok := expectedModules[m.Name]; ok {
			expectedModules[m.Name] = true
		}
	}

	for name, found := range expectedModules {
		if !found {
			t.Errorf("expected to find module %q but it was not loaded", name)
		}
	}
	t.Logf("Found %d modules", len(modules))
	for _, m := range modules {
		t.Logf(" -%s (path: %s)", m.Name, m.DirPath)
	}
}

func TestLoadAllSkipsHiddenDirs(t *testing.T) {
	// Créer un dossier temporaire avec une structure
	tmpDir := t.TempDir()

	// Créer un module normal
	normalDir := filepath.Join(tmpDir, "normal")
	os.MkdirAll(normalDir, 0755)
	writePKGBUILD(t, normalDir, "normal-mod")

	// Créer un module dans un dossier caché (doit être ignoré)
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	os.MkdirAll(hiddenDir, 0755)
	writePKGBUILD(t, hiddenDir, "hidden-mod")

	loader := NewBashLoader()
	modules, err := loader.LoadAll(tmpDir)
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(modules) != 1 {
		t.Errorf("expected 1 module, got %d", len(modules))
	}

	if len(modules) > 0 && modules[0].Name != "normal-mod" {
		t.Errorf("expected module 'normal-mod', got %q", modules[0].Name)
	}
}

func TestLoadAllSkipsBuildDirs(t *testing.T) {
	tmpDir := t.TempDir()

	// Module normal
	normalDir := filepath.Join(tmpDir, "mylib")
	os.MkdirAll(normalDir, 0755)
	writePKGBUILD(t, normalDir, "mylib")

	// Module dans bin/ (doit être ignoré)
	binDir := filepath.Join(tmpDir, "bin", "something")
	os.MkdirAll(binDir, 0755)
	writePKGBUILD(t, binDir, "bin-mod")

	// Module dans build/ (doit être ignoré)
	buildDir := filepath.Join(tmpDir, "build", "something")
	os.MkdirAll(buildDir, 0755)
	writePKGBUILD(t, buildDir, "build-mod")

	loader := NewBashLoader()
	modules, err := loader.LoadAll(tmpDir)
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(modules) != 1 {
		t.Errorf("expected 1 module, got %d", len(modules))
		for _, m := range modules {
			t.Logf("  found: %s", m.Name)
		}
	}
}

func writePKGBUILD(t *testing.T, dir, name string) {
	t.Helper()
	content := `pkgname=` + name + `
pkgdesc="Test module"
depends=()
makedepends=(clang)
source=(dummy.c)

produces() {
    echo "lib/lib` + name + `.a"
}

build() {
    echo "building..."
}
`
	path := filepath.Join(dir, "PKGBUILD")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write PKGBUILD: %v", err)
	}

	// Créer aussi dummy.c pour que la validation passe
	dummyPath := filepath.Join(dir, "dummy.c")
	if err := os.WriteFile(dummyPath, []byte("// dummy\n"), 0644); err != nil {
		t.Fatalf("failed to write dummy.c: %v", err)
	}
}
