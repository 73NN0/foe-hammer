package commands

import (
	"flag"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"os"
	"os/exec"
)


// Interface pour la config
type ConfigRepository interface {
	Get(key string) (string, error)
	GetAll() map[string]string
}

type BuildCommand struct {
	fs     *flag.FlagSet
	config  ConfigRepository
}

func NewBuildCommand(config  ConfigRepository) *BuildCommand {
	return &BuildCommand{
		fs:     flag.NewFlagSet("build", flag.ExitOnError),
		config: config,
	}
}

func (b *BuildCommand) Name() string           { return "build" }
func (b *BuildCommand) Description() string    { return "build command in fonction of target platform" }
func (c *BuildCommand) FlagSet() *flag.FlagSet { return c.fs }

func cflags(config  ConfigRepository, buildType string) string {

	var flags string

	if common, err := config.Get("cflagsCommon"); err != nil {
		return ""
	} else {
		flags += common
	}

	switch buildType {
	case "debug":
		if debug, err := config.Get("cflagsDebug"); err != nil {
			return flags
		} else {
			flags += debug
		}

	case "release":

		if release, err := config.Get("cflagsRelease"); err != nil {
			return flags
		} else {
			flags += release
		}
	default:
		return flags
	}

	return flags
}

func includeFlags(config  ConfigRepository, platform string) (flags string, err error) {
	rootDir, err := config.Get("rootDir")
	if err != nil {
		err = fmt.Errorf("No rootDir")
		flags = ""
		return
	}

	flags = "-I" + platform + "-I" + rootDir + "/include/" + "-I/include/sys"
	return
}

func executeCommand(commands [][]string, errorType error, cmdDir string) error {
	if len(commands) == 0 {
		return fmt.Errorf("no commands provided")
	}

	for _, cmd := range commands {
		command := exec.Command(cmd[0], cmd[1:]...)
		command.Stdout = os.Stdout
		command.Stderr = nil
		command.Env = append(os.Environ(), "LANG=C")
		if len(cmdDir) > 0 {
			command.Dir = cmdDir
		}
		if err := command.Run(); err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				return fmt.Errorf("%w: %v", errorType, err)
			}
			return err
		}
	}

	return nil
}

// discoverLibs trouve tous les répertoires lib* sous ROOT_DIR/src
func discoverLibs(rootDir string) ([]string, error) {
	baseDir := filepath.Join(rootDir, "src")

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", baseDir, err)
	}

	var libDirs []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "lib") {
			libDirs = append(libDirs, filepath.Join(baseDir, entry.Name()))
		}
	}

	return libDirs, nil
}

// sourcesLib trouve tous les fichiers .c dans un répertoire
func sourcesLib(libDir string) ([]string, error) {
	var sources []string

	err := filepath.Walk(libDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".c" {
			sources = append(sources, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk %s: %w", libDir, err)
	}

	return sources, nil
}

// compileLibSources - version optimisée
func compileLibSources(config ConfigRepository, libName string, sources []string, objDir string) ([]string, error) {
	outDir := filepath.Join(objDir, libName)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create obj dir: %w", err)
	}

	objects := make([]string, 0, len(sources))

	// Compiler un par un au lieu de batcher
	for _, src := range sources {
		objName := strings.TrimSuffix(filepath.Base(src), ".c") + ".o"
		objPath := filepath.Join(outDir, objName)
		objects = append(objects, objPath)

		// Construire et exécuter immédiatement
		cmd := buildCompileCommand(config, src, objPath)

		fmt.Printf("  CC %s\n", src)

		if err := executeCommand([][]string{cmd}, fmt.Errorf("compilation failed"), ""); err != nil {
			return nil, fmt.Errorf("failed to compile %s: %w", src, err)
		}
	}

	return objects, nil
}

func buildCompileCommand(config ConfigRepository, src, obj string) []string {
	var cc string
	var ccStd string
	var CFlags string

	if v, err := config.Get("cc"); err == nil {
		cc = v
	}

	if v, err := config.Get("ccStd"); err == nil {
		ccStd = v
	}

	if v, err := config.Get("cFlags"); err == nil {
		CFlags = v
	}


	inFlags, err := includeFlags(config, platform)

	flags := cflags(config, buildType)

	cmd := []string{cc, ccStd}
	cmd = append(cmd, strings.Fields(CFlags)...)
	cmd = append(cmd, strings.Fields(config.IncludeFlags)...)
	cmd = append(cmd, "-c", src, "-o", obj)
	return cmd
}

// buildLib crée une archive statique .a à partir des .o
func buildLib(libName string, objects []string, libOutputDir string) error {
	if len(objects) == 0 {
		fmt.Printf("  → no objects for %s, skipping\n", libName)
		return nil
	}

	if err := os.MkdirAll(libOutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create lib dir: %w", err)
	}

	outLib := filepath.Join(libOutputDir, libName+".a")

	// Commande ar pour créer l'archive
	cmd := []string{"ar", "rcs", outLib}
	cmd = append(cmd, objects...)

	commands := [][]string{cmd}

	if err := executeCommand(commands, fmt.Errorf("archiving failed"), ""); err != nil {
		return err
	}

	fmt.Printf("  ✓ %s built\n", outLib)
	return nil
}

// buildAllLibs orchestre la compilation de toutes les libs
func buildAllLibs(config *Config, rootDir, objDir, libOutputDir string) error {
	// 1. Découvrir les libs
	libDirs, err := discoverLibs(rootDir)
	if err != nil {
		return err
	}

	if len(libDirs) == 0 {
		fmt.Printf("[build_all_libs] No libraries found under %s/src\n", rootDir)
		return nil
	}

	fmt.Printf("Found %d libraries\n", len(libDirs))

	// 2. Pour chaque lib
	for _, libDir := range libDirs {
		libName := filepath.Base(libDir)
		fmt.Printf("\n=== Building library: %s ===\n", libName)

		// 2.1 Trouver les sources
		sources, err := sourcesLib(libDir)
		if err != nil {
			return fmt.Errorf("failed to find sources for %s: %w", libName, err)
		}

		if len(sources) == 0 {
			fmt.Printf("  → no sources in %s\n", libName)
			continue
		}

		// 2.2 Compiler les sources
		objects, err := compileLibSources(config, libName, sources, objDir)
		if err != nil {
			return fmt.Errorf("failed to compile %s: %w", libName, err)
		}

		// 2.3 Créer l'archive
		if err := buildLib(libName, objects, libOutputDir); err != nil {
			return fmt.Errorf("failed to build lib %s: %w", libName, err)
		}
	}

	return nil
}
func (b *BuildCommand) Run(args []string) error {
	var _argsLen int
	var buildType string
	var platform string
	var output string

	_argsLen = len(args)
	if _argsLen < 1 {
		// TODO error.Is
		return errors.New("No args")
	}

	b.fs.StringVar(&buildType, "type", "release", "define build type [release|target] by default : release")
	b.fs.StringVar(&platform, "platform", runtime.GOOS, "define target build platform")
	b.fs.StringVar(&output, "output", "chief", "define output name")
	// TODO runtime.GOARCH

	b.fs.Parse(args)


	if err != nil {
		return err
	}

	return nil
}
