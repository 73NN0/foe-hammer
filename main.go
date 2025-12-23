package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"

	commands "github.com/73NN0/foe-hammer"
)

// Interface pour la config
type ConfigRepository interface {
	Get(key string) (string, error)
	GetAll() map[string]string
}

// Implémentation simple
type InMemoryConfig struct {
	data map[string]string
}

func NewInMemoryConfig() *InMemoryConfig {
	return &InMemoryConfig{
		data: map[string]string{
			// todo improve this
			"rootDir":  "sys",
			"buildDir": "bin",
			"cc":       "clang",
			"cflagsCommon": `-Wall -Wextra -Wpedantic \
		-Wshadow -Wcast-align -Wunused -Wold-style-definition \
		-Wmissing-prototypes -Wno-unused-parameter -Werror \
		-Wstrict-prototypes -Wpointer-arith -Wwrite-strings \
		-Wconversion -Wformat=2 -Wformat-security \
		-Wunreachable-code -Wundef -Wbad-function-cast \
		-Wdouble-promotion -Wmissing-include-dirs \
		-Winit-self -Wmissing-noreturn -fno-common \
		-fstack-protector-strong`,

			"cflagsRelease": `-O2 -DNDEBUG -DDEBUG_MEMORY=0 -fomit-frame-pointer -march=native -D_FORTIFY_SOURCE=2`,

			"cflagsDebug": `-g3 -O0 -DDEBUG -DDEBUG_MEMORY=1 -ftrapv`,
			"linkFlags" : "-lSDL2 -lm",
		},
	}
}

func (c *InMemoryConfig) Get(key string) (string, error) {
	if val, ok := c.data[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key not found: %s", key)
}

func (c *InMemoryConfig) GetAll() map[string]string {
	return c.data
}

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

func (r *Registry) List() []Command {
	cmds := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// ============================================================================
// Help Command - spécial
// ============================================================================

type HelpCommand struct {
	fs       *flag.FlagSet
	registry *Registry
}

func NewHelpCommand(registry *Registry) *HelpCommand {
	return &HelpCommand{
		fs:       flag.NewFlagSet("help", flag.ExitOnError),
		registry: registry,
	}
}

func (c *HelpCommand) Name() string           { return "help" }
func (c *HelpCommand) Description() string    { return "Show help information" }
func (c *HelpCommand) FlagSet() *flag.FlagSet { return c.fs }

func (c *HelpCommand) Run(args []string) error {
	if len(args) > 0 {
		// Help pour une commande spécifique
		cmdName := args[0]
		if cmd, ok := c.registry.Get(cmdName); ok {
			fmt.Printf("Usage: foe %s [options]\n\n", cmd.Name())
			fmt.Printf("%s\n\n", cmd.Description())
			fmt.Println("Options:")
			cmd.FlagSet().PrintDefaults()
			return nil
		}
		return fmt.Errorf("unknown command: %s", cmdName)
	}

	subtitles := []string{
		// Halo
		"Echo 419 inbound!", // Foe Hammer's callsign
		"This is Foe Hammer, arriving at LZ",
		"Somebody order a Warthog?",
		"Finish the fight",
		"Don't make a girl a promise you can't keep",

		// DRG
		"ROCK AND STONE!",
		"FOR KARL!",
		"If you don't Rock and Stone, you ain't comin' home!",

		// DOOM
		"RIP AND TEAR... UNTIL THE BUILD IS DONE",
		"KNEE-DEEP IN COMPILATION ERRORS",

		// Warhammer
		"For the emperor",
		"purge heretics",

		// Minecraft
		"Also try CMake!", // ironic

		// Misc
		"because I'm too lazy to learn cmake",
		"Praise be to Space king",
		"no girls allowed",
		"activate plan 9",
	}

	dice := rand.Intn(len(subtitles))

	name := `
███████╗ ██████╗ ███████╗    ██╗  ██╗ █████╗ ███╗   ███╗███╗   ███╗███████╗██████╗ 
██╔════╝██╔═══██╗██╔════╝    ██║  ██║██╔══██╗████╗ ████║████╗ ████║██╔════╝██╔══██╗
█████╗  ██║   ██║█████╗      ███████║███████║██╔████╔██║██╔████╔██║█████╗  ██████╔╝
██╔══╝  ██║   ██║██╔══╝      ██╔══██║██╔══██║██║╚██╔╝██║██║╚██╔╝██║██╔══╝  ██╔══██╗
██║     ╚██████╔╝███████╗    ██║  ██║██║  ██║██║ ╚═╝ ██║██║ ╚═╝ ██║███████╗██║  ██║
╚═╝      ╚═════╝ ╚══════╝    ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝                                                                        
 `

	fmt.Printf("%s \n %s \n", name, subtitles[dice])
	fmt.Println("\nUsage: foe <command> [options]")
	fmt.Println("\nAvailable commands:")

	cmds := c.registry.List()
	maxLen := 0
	for _, cmd := range cmds {
		if len(cmd.Name()) > maxLen {
			maxLen = len(cmd.Name())
		}
	}

	for _, cmd := range cmds {
		padding := strings.Repeat(" ", maxLen-len(cmd.Name())+2)
		fmt.Printf("  %s%s%s\n", cmd.Name(), padding, cmd.Description())
	}

	fmt.Println("\nUse 'foe help <command>' for more information about a command.")
	return nil
}


// CLI

type CLI struct {
	registry *Registry
	config   ConfigRepository
}

func NewCLI(config ConfigRepository) *CLI {
	cli := &CLI{
		registry: NewRegistry(),
		config:   config,
	}
	cli.registry.Register(NewHelpCommand(cli.registry))
	cli.registry.Register(commands.NewBuildCommand(config))
	return cli
}

func (c *CLI) Run(args []string) error {
	if len(args) < 1 {
		// no command -> display help
		fmt.Println("no")

		return c.registry.commands["help"].Run([]string{})
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	cmd, ok := c.registry.Get(cmdName)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error : unknown command '%s'\n", cmdName)
		c.registry.commands["help"].Run([]string{})
		return fmt.Errorf("unknown command: %s", cmdName)

	}

	return cmd.Run(cmdArgs)
}

func main() {
	cli := NewCLI(NewInMemoryConfig())

	if err := cli.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
