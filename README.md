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

*foe-hammer injects [environment variables](adapters/context/readme.md) into your hooks*

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