package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
)

const (
	version string = `poc version (V0.0.0)`
	name    string = `
███████╗ ██████╗ ███████╗    ██╗  ██╗ █████╗ ███╗   ███╗███╗   ███╗███████╗██████╗
██╔════╝██╔═══██╗██╔════╝    ██║  ██║██╔══██╗████╗ ████║████╗ ████║██╔════╝██╔══██╗
█████╗  ██║   ██║█████╗      ███████║███████║██╔████╔██║██╔████╔██║█████╗  ██████╔╝
██╔══╝  ██║   ██║██╔══╝      ██╔══██║██╔══██║██║╚██╔╝██║██║╚██╔╝██║██╔══╝  ██╔══██╗
██║     ╚██████╔╝███████╗    ██║  ██║██║  ██║██║ ╚═╝ ██║██║ ╚═╝ ██║███████╗██║  ██║
╚═╝      ╚═════╝ ╚══════╝    ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝
 `
)

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

	fmt.Printf("%s \n %s \t\t\t%s \n", name, version, subtitles[dice])
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
