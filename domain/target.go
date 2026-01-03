package domain

import (
	"fmt"
	"runtime"
)

// for what platform we want to compile to
type Target struct {
	OS   string // win32 linux darwin android
	Arch string // x86_64, arm64 etc...
}

func (t Target) String() string {
	return fmt.Sprintf("%v-%v", t.Arch, t.OS)
}

func NewTarget() Target {
	target := Target{}
	target.OS = runtime.GOOS
	target.Arch = runtime.GOARCH
	return target
}
