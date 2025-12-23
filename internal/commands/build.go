package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ConfigRepository defines the interface for configuration access
type ConfigRepository interface {
	Get(key string) (string, error)
	GetAll() map[string]string
}

type BuildCommand struct {
	fs     *flag.FlagSet
	config ConfigRepository
}

func NewBuildCommand(config ConfigRepository) *BuildCommand {
	return &BuildCommand{
		fs:     flag.NewFlagSet("build", flag.ExitOnError),
		config: config,
	}
}

func (b *BuildCommand) Name() string           { return "build" }
func (b *BuildCommand) Description() string    { return "Build project for target platform" }
func (b *BuildCommand) FlagSet() *flag.FlagSet { return b.fs }

// BuildContext holds all build parameters parsed from flags
type BuildContext struct {
	BuildType string
	Platform  string
	Output    string
	RootDir   string
	BuildDir  string
	ObjDir    string
	LibDir    string
}

func (b *BuildCommand) Run(args []string) error {
	if len(args) < 1 {
		return errors.New("no arguments provided")
	}

	var buildType, platform, output string

	b.fs.StringVar(&buildType, "type", "release", "build type [release|debug]")
	b.fs.StringVar(&platform, "platform", runtime.GOOS, "target platform")
	b.fs.StringVar(&output, "output", "chief", "output binary name")

	if err := b.fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	rootDir, err := b.config.Get("rootDir")
	if err != nil {
		return fmt.Errorf("rootDir not configured: %w", err)
	}

	buildDir, err := b.config.Get("buildDir")
	if err != nil {
		buildDir = "bin"
	}

	ctx := &BuildContext{
		BuildType: buildType,
		Platform:  platform,
		Output:    output,
		RootDir:   rootDir,
		BuildDir:  buildDir,
		ObjDir:    filepath.Join(buildDir, buildType, platform, "obj"),
		LibDir:    filepath.Join(buildDir, buildType, platform, "lib"),
	}

	return b.buildAll(ctx)
}

func (b *BuildCommand) buildAll(ctx *BuildContext) error {
	libDirs, err := discoverLibs(ctx.RootDir)
	if err != nil {
		return err
	}

	if len(libDirs) == 0 {
		fmt.Printf("[build] No libraries found under %s/src\n", ctx.RootDir)
		return nil
	}

	fmt.Printf("Found %d libraries\n", len(libDirs))

	for _, libDir := range libDirs {
		libName := filepath.Base(libDir)
		fmt.Printf("\n=== Building library: %s ===\n", libName)

		sources, err := sourcesLib(libDir)
		if err != nil {
			return fmt.Errorf("failed to find sources for %s: %w", libName, err)
		}

		if len(sources) == 0 {
			fmt.Printf("  → no sources in %s\n", libName)
			continue
		}

		objects, err := b.compileLibSources(ctx, libName, sources)
		if err != nil {
			return fmt.Errorf("failed to compile %s: %w", libName, err)
		}

		if err := buildLib(libName, objects, ctx.LibDir); err != nil {
			return fmt.Errorf("failed to build lib %s: %w", libName, err)
		}
	}

	return nil
}

func (b *BuildCommand) cflags(buildType string) string {
	var flags string

	if common, err := b.config.Get("cflagsCommon"); err == nil {
		flags = common
	}

	switch buildType {
	case "debug":
		if debug, err := b.config.Get("cflagsDebug"); err == nil {
			flags += " " + debug
		}
	case "release":
		if release, err := b.config.Get("cflagsRelease"); err == nil {
			flags += " " + release
		}
	}

	return flags
}

func (b *BuildCommand) includeFlags(ctx *BuildContext) string {
	var parts []string

	if ctx.Platform != "" {
		parts = append(parts, "-I"+ctx.Platform)
	}

	parts = append(parts, "-I"+ctx.RootDir+"/include")
	parts = append(parts, "-I"+ctx.RootDir+"/include/sys")

	return strings.Join(parts, " ")
}

func (b *BuildCommand) buildCompileCommand(ctx *BuildContext, src, obj string) []string {
	cc := "cc"
	ccStd := "-std=c11"

	if v, err := b.config.Get("cc"); err == nil {
		cc = v
	}
	if v, err := b.config.Get("ccStd"); err == nil {
		ccStd = v
	}

	cflags := b.cflags(ctx.BuildType)
	includeFlags := b.includeFlags(ctx)

	cmd := []string{cc, ccStd}
	cmd = append(cmd, strings.Fields(cflags)...)
	cmd = append(cmd, strings.Fields(includeFlags)...)
	cmd = append(cmd, "-c", src, "-o", obj)

	return cmd
}

func (b *BuildCommand) compileLibSources(ctx *BuildContext, libName string, sources []string) ([]string, error) {
	outDir := filepath.Join(ctx.ObjDir, libName)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create obj dir: %w", err)
	}

	objects := make([]string, 0, len(sources))

	for _, src := range sources {
		objName := strings.TrimSuffix(filepath.Base(src), ".c") + ".o"
		objPath := filepath.Join(outDir, objName)
		objects = append(objects, objPath)

		cmd := b.buildCompileCommand(ctx, src, objPath)

		fmt.Printf("  CC %s\n", src)

		if err := executeCommand([][]string{cmd}, ""); err != nil {
			return nil, fmt.Errorf("failed to compile %s: %w", src, err)
		}
	}

	return objects, nil
}

func executeCommand(commands [][]string, cmdDir string) error {
	if len(commands) == 0 {
		return fmt.Errorf("no commands provided")
	}

	for _, cmd := range commands {
		command := exec.Command(cmd[0], cmd[1:]...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Env = append(os.Environ(), "LANG=C")

		if cmdDir != "" {
			command.Dir = cmdDir
		}

		if err := command.Run(); err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				return fmt.Errorf("command failed: %w", err)
			}
			return err
		}
	}

	return nil
}

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

func buildLib(libName string, objects []string, libOutputDir string) error {
	if len(objects) == 0 {
		fmt.Printf("  → no objects for %s, skipping\n", libName)
		return nil
	}

	if err := os.MkdirAll(libOutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create lib dir: %w", err)
	}

	outLib := filepath.Join(libOutputDir, libName+".a")

	cmd := []string{"ar", "rcs", outLib}
	cmd = append(cmd, objects...)

	if err := executeCommand([][]string{cmd}, ""); err != nil {
		return err
	}

	fmt.Printf("  ✓ %s built\n", outLib)
	return nil
}
