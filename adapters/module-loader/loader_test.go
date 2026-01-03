package moduleloader

import (
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
