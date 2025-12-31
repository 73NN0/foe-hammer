# foe-hammer

> A build orchestrator that doesn't care how you build, just what you produce

## What it does

1. Scans for modules (PKGBUILD files)
2. Builds a dependency graph
3. Computes build order (topological sort)
4. Executes each module's `build()` hook in order

**foe-hammer does NOT compile anything.** It just orchestrates.

## PKGBUILD format

```bash
pkgname=mylib
pkgdesc="My library"
produces=(lib/libmylib.a)
depends=(otherliba otherlibb)
makedepends=(clang)
source=(foo.c bar.c)

build() {
    mkdir -p "$FOE_OBJDIR" "$FOE_LIBDIR"
    clang -c "$FOE_SRCDIR/foo.c" -o "$FOE_OBJDIR/foo.o"
    clang -c "$FOE_SRCDIR/bar.c" -o "$FOE_OBJDIR/bar.o"
    ar rcs "$FOE_LIBDIR/libmylib.a" "$FOE_OBJDIR"/*.o
}
```

## Environment variables

foe-hammer injects these into your `build()` hook:

| Variable | Description |
|----------|-------------|
| `FOE_HOST_OS` | Host OS (darwin, linux...) |
| `FOE_HOST_ARCH` | Host arch (amd64, arm64...) |
| `FOE_TARGET_OS` | Target OS |
| `FOE_TARGET_ARCH` | Target arch |
| `FOE_OUTDIR` | Build output root |
| `FOE_LIBDIR` | Where to put .a files |
| `FOE_BINDIR` | Where to put executables |
| `FOE_OBJDIR` | Where to put .o files (per module) |
| `FOE_SRCDIR` | Module source directory |
| `FOE_MODULE_NAME` | Current module name |

## Architecture

```
PKGBUILD → ModuleLoader → Module
                            ↓
                       BuildGraph → order (topo sort)
                            ↓
                       Orchestrator
                            ↓
                       ContextProvider → env vars
                            ↓
                       HookRunner → executes build()
                            ↓
                       Executor
```

## Usage

```bash
go build -o foe ./main.go
./foe ./myproject
```