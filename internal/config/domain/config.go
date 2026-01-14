package domain

type ProjectConfig struct {
	ID               int
	RootDir          string
	ManifestFilename string
	IgnoreDirs       []string
	OutDirDefault    string
}
