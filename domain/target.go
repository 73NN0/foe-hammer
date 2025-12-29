package domain

import (
	"fmt"
)

// for what platform we want to compile to
type Target struct {
	OS   string // win32 linux darwin android
	Arch string // x86_64, arm64 etc...
}

func (t Target) String() string {
	return fmt.Sprintf("%v-%v", t.Arch, t.OS)
}

// from which platform we compile
type Host struct {
	OS   string
	Arch string
}

func (h Host) String() string {
	return fmt.Sprintf("%v-%v", h.Arch, h.OS)
}

func (h Host) CrossCompilingTo(t Target) bool {
	return h.OS != t.OS || h.Arch != t.Arch
}
