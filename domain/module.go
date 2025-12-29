package domain

// represent a parsed PKGBUILD

type Module struct {
	Name        string
	DirPath     string // Directory where the module lives
	Path        string // full path (abs )
	Description string
	Produces    []string // relatifs paths of build artefacts
	Depends     []string // dependency modules
	MakeDepends []string // external dependency ( SDL2 etc...)
	Sources     []string // sources files
	Strategy    Strategy // system or custom
	BuildFunc   string   // build hook content if custom stragey choosen
}

type Strategy string

const (
	StrategySystem Strategy = "system"
	StrategyCustom Strategy = "custom"
)
