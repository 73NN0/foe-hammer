package cli

import (
	"fmt"
	"os"

	commands "github.com/73NN0/foe-hammer/internal/commands"
	config "github.com/73NN0/foe-hammer/internal/config/domain"
	registry "github.com/73NN0/foe-hammer/internal/registry/domain"
)

type CLI struct {
	registry *registry.Registry
	config   config.ConfigRepository
}

func NewCLI(config config.ConfigRepository) *CLI {
	cli := &CLI{
		registry: registry.NewRegistry(),
		config:   config,
	}
	cli.registry.Register(commands.NewHelpCommand(cli.registry))
	cli.registry.Register(commands.NewBuildCommand(config))
	return cli
}

func (c *CLI) Run(args []string) error {
	if len(args) < 1 {
		// no command -> display help
		fmt.Println("no")

		return c.registry.RunCommand("help", []string{})
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	cmd, ok := c.registry.Get(cmdName)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error : unknown command '%s'\n", cmdName)
		c.registry.RunCommand("help", []string{})
		return fmt.Errorf("unknown command: %s", cmdName)

	}

	return cmd.Run(cmdArgs)
}
