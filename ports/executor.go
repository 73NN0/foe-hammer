package ports

// Executor knows how execute shell commands
// TODO : Question only shell ?
// Repository pattern and one implementation for shell commands ?
// that allow to move aways from shell
// maybe after poc version

type Executor interface {
	Run(cmd string, args []string, workDir string) error
	RunScript(script string, workDir string) error
}
