package domain

// represent a parsed PKGBUILD

type Module struct {
	Name        string
	Description string
	Produces    []string // relatifs paths of build artefacts
	Depends     []string // dependency modules
	MakeDepends []string // external dependency ( SDL2 etc...)
	Sources     []string // sources files
	Strategy    strategy // system or custom
	BuildFunc   string   // build hook content if custom stragey choosen
}

type strategy string

const (
	StrategySystem strategy = "system"
	StrategyCustom strategy = "custom"
)
