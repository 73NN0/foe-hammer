package domain

// for what platform we want to compile to
type Target struct {
	OS   string // window linux darwin android
	Arch string // x86_64, arm64 etc...
}

// from which platform we compile
type Host struct {
	OS   string
	Arch string
}

func (h Host) CrossCompilingTo(t Target) bool {
	return h.OS != t.OS || h.Arch != t.Arch
}
