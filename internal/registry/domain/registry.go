package cli

import (
	"flag"
	"fmt"
)

type Command interface {
	Name() string
	Description() string
	Run( /*args*/ []string) error
	FlagSet() *flag.FlagSet
}

// Registry
type Registry struct {
	commands map[string]Command
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]Command),
	}
}

func (r *Registry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

func (r *Registry) Get(name string) (Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

func (r *Registry) RunCommand(name string, args []string) error {
	if name == "" {
		return fmt.Errorf("no command name provided")
	}
	return r.commands["help"].Run(args)
}

func (r *Registry) List() []Command {
	cmds := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}
