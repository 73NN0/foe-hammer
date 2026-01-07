package ports

import "io"

type Executor interface {
	Run(cmd string, args []string, workDir string, stdout, stderr io.Writer) error
}
