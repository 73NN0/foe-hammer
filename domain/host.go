package domain

import (
	"fmt"
	"runtime"
)

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

func NewHost() Host {
	host := Host{}
	host.OS = runtime.GOOS
	host.Arch = runtime.GOARCH
	return host
}
