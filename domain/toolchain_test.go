package domain

import "testing"

var happyPath = []CrossRule{
	// Natif = simple

	// Cross vers Windows depuis n'importe où
	{HostOS: "darwin", TargetOS: "windows", Arch: "x86_64", Config: CrossConfig{
		Compiler: "clang",
		Flags:    []string{"--target=x86_64-pc-windows-msvc"},
		Linker:   "lld-link",
	}},
}

func TestMatch(t *testing.T) {
	tests := []struct {
		name      string
		CrossRule CrossRule
		host      Host
		target    Target
		want      bool
	}{
		{
			name: "happy path without *",
			CrossRule: CrossRule{HostOS: "linux", TargetOS: "linux", Arch: "x86_64", Config: CrossConfig{
				Compiler: "clang",
			}},
			host:   Host{OS: "linux", Arch: "x86_64"},
			target: Target{OS: "linux", Arch: "x86_64"},
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			get := matches(tt.CrossRule, tt.host, tt.target)

			if get != tt.want {
				t.Errorf("got %t, want %t", get, tt.want)
			}
		})
	}

}
