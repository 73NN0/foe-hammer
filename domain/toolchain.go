package domain

type CrossRule struct {
	HostOS   string // "*" = any
	TargetOS string
	Arch     string // "*" = any
	Config   CrossConfig
}

type CrossConfig struct {
	Compiler    string   // "clang", "zig"
	Flags       []string // flags additionnels
	Linker      string   // "lld", "link.exe"
	NeedsPrefix string   // "wine" si on doit exécuter via wine
}

// Clé = "host-os:target-os:target-arch"
var CrossRules = []CrossRule{
	// Natif = simple
	{HostOS: "*", TargetOS: "*", Arch: "*", Config: CrossConfig{
		Compiler: "clang",
		// host == target, pas de flags spéciaux
	}},

	// Cross vers Windows depuis n'importe où
	{HostOS: "*", TargetOS: "windows", Arch: "x86_64", Config: CrossConfig{
		Compiler: "clang",
		Flags:    []string{"--target=x86_64-pc-windows-msvc"},
		Linker:   "lld-link",
	}},
}

func matches(rule CrossRule, host Host, target Target) bool {
	hostOSMatch := rule.HostOS == "*" || rule.HostOS == host.OS
	targetOSMatch := rule.TargetOS == "*" || rule.TargetOS == target.OS
	archMatch := rule.Arch == "*" || rule.Arch == target.Arch

	return hostOSMatch && targetOSMatch && archMatch
}

func GetCrossConfig(host Host, target Target) CrossConfig {
	// Parcourir les règles, la plus spécifique gagne <= Thinking Do I need to iterate ?
	for _, rule := range CrossRules {
		if matches(rule, host, target) {
			return rule.Config
		}
	}
	// Fallback natif
	return CrossConfig{Compiler: "clang"}
}
