package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type LinkCommand struct {
	fs     *flag.FlagSet
	config ConfigRepository
}

func NewLinkCommand(config ConfigRepository) *LinkCommand {
	return &LinkCommand{
		fs:     flag.NewFlagSet("link", flag.ExitOnError),
		config: config,
	}
}

func (l *LinkCommand) Name() string           { return "link" }
func (l *LinkCommand) Description() string    { return "Link compiled objects into final executable" }
func (l *LinkCommand) FlagSet() *flag.FlagSet { return l.fs }

// LinkContext holds all link parameters
type LinkContext struct {
	BuildType string
	Platform  string
	Output    string
	RootDir   string
	BuildDir  string
	ObjDir    string
	LibDir    string
	LibOrder  []string
	EntryFile string
}

func (l *LinkCommand) Run(args []string) error {
	var buildType, platform, output, libOrder, entryFile string

	l.fs.StringVar(&buildType, "type", "release", "build type [release|debug]")
	l.fs.StringVar(&platform, "platform", runtime.GOOS, "target platform")
	l.fs.StringVar(&output, "output", "chief", "output binary name")
	l.fs.StringVar(&libOrder, "lib-order", "", "library link order (comma-separated)")
	l.fs.StringVar(&entryFile, "entry", "", "entry point source file")

	if err := l.fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	rootDir, err := l.config.Get("rootDir")
	if err != nil {
		return fmt.Errorf("rootDir not configured: %w", err)
	}

	buildDir, err := l.config.Get("buildDir")
	if err != nil {
		buildDir = "bin"
	}

	// Default entry file if not specified
	if entryFile == "" {
		entryFile = filepath.Join(rootDir, "src", "cmd", "engine", "engine.c")
	}

	// Parse lib order
	var libs []string
	if libOrder != "" {
		libs = strings.Split(libOrder, ",")
		for i := range libs {
			libs[i] = strings.TrimSpace(libs[i])
		}
	}

	basePath := filepath.Join(buildDir, buildType, platform)

	ctx := &LinkContext{
		BuildType: buildType,
		Platform:  platform,
		Output:    output,
		RootDir:   rootDir,
		BuildDir:  buildDir,
		ObjDir:    filepath.Join(basePath, "obj"),
		LibDir:    filepath.Join(basePath, "lib"),
		LibOrder:  libs,
		EntryFile: entryFile,
	}

	return l.link(ctx)
}

func (l *LinkCommand) link(ctx *LinkContext) error {
	exePath := filepath.Join(ctx.BuildDir, ctx.BuildType, ctx.Platform, ctx.Output)

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(exePath), 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	fmt.Println("=== Linking final executable ===")

	// Discover static libraries
	libArgs, err := l.discoverArchives(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("[link] Using libraries: %s\n", strings.Join(libArgs, " "))

	// Find platform objects
	platformObjs, err := l.discoverPlatformObjects(ctx)
	if err != nil {
		return err
	}

	// Build link command
	cmd := l.buildLinkCommand(ctx, exePath, libArgs, platformObjs)

	fmt.Printf("[link] %s\n", strings.Join(cmd, " "))

	if err := executeCommand([][]string{cmd}, ""); err != nil {
		return fmt.Errorf("linking failed: %w", err)
	}

	fmt.Printf("✓ Build OK → %s\n", exePath)
	return nil
}

func (l *LinkCommand) discoverArchives(ctx *LinkContext) ([]string, error) {
	var args []string

	// Add library search path
	args = append(args, "-L"+ctx.LibDir)

	if len(ctx.LibOrder) > 0 {
		// Follow specified order
		for _, name := range ctx.LibOrder {
			libPath := filepath.Join(ctx.LibDir, "lib"+name+".a")
			if _, err := os.Stat(libPath); os.IsNotExist(err) {
				fmt.Printf("[warn] expected lib%s.a not found in %s\n", name, ctx.LibDir)
				continue
			}
			args = append(args, "-l"+name)
		}
	} else {
		// Auto-discover libraries
		entries, err := os.ReadDir(ctx.LibDir)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("[warn] lib directory %s does not exist\n", ctx.LibDir)
				return args, nil
			}
			return nil, fmt.Errorf("failed to read lib dir: %w", err)
		}

		foundAny := false
		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, "lib") && strings.HasSuffix(name, ".a") {
				// Extract library name: libfoo.a -> foo
				libName := strings.TrimPrefix(name, "lib")
				libName = strings.TrimSuffix(libName, ".a")
				args = append(args, "-l"+libName)
				foundAny = true
			}
		}

		if !foundAny {
			fmt.Printf("[warn] no static libraries (.a) found in %s\n", ctx.LibDir)
		}
	}

	return args, nil
}

func (l *LinkCommand) discoverPlatformObjects(ctx *LinkContext) ([]string, error) {
	platformObjDir := filepath.Join(ctx.ObjDir, ctx.Platform)

	if _, err := os.Stat(platformObjDir); os.IsNotExist(err) {
		return nil, nil
	}

	var objects []string

	entries, err := os.ReadDir(platformObjDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read platform obj dir: %w", err)
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".o") {
			objects = append(objects, filepath.Join(platformObjDir, entry.Name()))
		}
	}

	return objects, nil
}

func (l *LinkCommand) buildLinkCommand(ctx *LinkContext, exePath string, libArgs, platformObjs []string) []string {
	cc := "cc"
	ccStd := "-std=c11"

	if v, err := l.config.Get("cc"); err == nil {
		cc = v
	}
	if v, err := l.config.Get("ccStd"); err == nil {
		ccStd = v
	}

	cflags := l.cflags(ctx.BuildType)
	includeFlags := l.includeFlags(ctx)
	linkFlags := l.linkFlags()

	cmd := []string{cc, ccStd}
	cmd = append(cmd, strings.Fields(cflags)...)
	cmd = append(cmd, strings.Fields(includeFlags)...)

	// Entry point source
	cmd = append(cmd, ctx.EntryFile)

	// Platform objects
	cmd = append(cmd, platformObjs...)

	// Library args (-L and -l flags)
	cmd = append(cmd, libArgs...)

	// System link flags
	cmd = append(cmd, strings.Fields(linkFlags)...)

	// Output
	cmd = append(cmd, "-o", exePath)

	return cmd
}

func (l *LinkCommand) cflags(buildType string) string {
	var flags string

	if common, err := l.config.Get("cflagsCommon"); err == nil {
		flags = common
	}

	switch buildType {
	case "debug":
		if debug, err := l.config.Get("cflagsDebug"); err == nil {
			flags += " " + debug
		}
	case "release":
		if release, err := l.config.Get("cflagsRelease"); err == nil {
			flags += " " + release
		}
	}

	return flags
}

func (l *LinkCommand) includeFlags(ctx *LinkContext) string {
	var parts []string

	if ctx.Platform != "" {
		parts = append(parts, "-I"+ctx.Platform)
	}

	parts = append(parts, "-I"+ctx.RootDir+"/include")
	parts = append(parts, "-I"+ctx.RootDir+"/include/sys")

	return strings.Join(parts, " ")
}

func (l *LinkCommand) linkFlags() string {
	if flags, err := l.config.Get("linkFlags"); err == nil {
		return flags
	}
	return "-lSDL2 -lm"
}
