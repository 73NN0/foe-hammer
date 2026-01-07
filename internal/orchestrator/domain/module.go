package domain

// represent a parsed PKGBUILD

type Module struct {
	Name        string
	DirPath     string // Directory where the module lives
	Path        string // full path (abs)
	Description string
	Produces    []string // relatifs paths of build artefacts
	Depends     []string // dependency modules
	MakeDepends []string // external dependency (SDL2 etc...)
	Sources     []string // sources files
}
