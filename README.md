<div align="left" style="position: relative;">
<img src="https://github.com/user-attachments/assets/9386181d-86bc-4093-a374-af509470ee18" align="right" width="40%" style="margin: -20px 0 0 20px; z-index: 3;">
<h1>FOE-HAMMER</h1>
<p align="left">
	<em><code>❯ A build orchestrator that doesn't care how you build, just what you produce </code></em>
</p>
<p align="left">
	<img src="https://img.shields.io/github/languages/top/73NN0/foe-hammer?style=default&color=05b926" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/73NN0/foe-hammer?style=default&color=05b926" alt="repo-language-count">
</p>
<p align="left"><!-- default option, no dependency badges. -->
    <strong>foe-hammer does NOT compile anything</strong> It just orchestrates.
</p>
<p align="left">
	Under active development 
</p>
</div>
<br clear="right">


## What it does

1. Scans for modules (PKGBUILD files)
2. Builds a dependency graph
3. Computes build order (topological sort)
4. Executes each module's `build()` hook in order


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
                       HookRunner → executes Produces()
                            ↓
                       HookRunner → executes build()
                            ↓
                       Executor
```

## Build

```bash
    make
```
