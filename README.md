<img width="768" height="512" alt="image" src="https://github.com/user-attachments/assets/f0cdbb2a-9e7e-4cf3-a0b8-6544c1186787" />

# foe-hammer

*Yes I'm a fan of the Halo franchise*

> A build orchestrator that doesn't care how you build, just what you produce

## why 

I'm just tired to convince friends to code with me in C projects because I'm not using the "right" tool ( aka Cmake )

I personaly prefer to write script and dealing with the compiler. 

So I'm trying to build a tool to concile both approach.

It's under heavely development

## The idea

You have a C project with multiple modules.
Some use CMake, some use hand-written Makefiles, some just shell scripts.

**foe-hammer** says: "Tell me what each module produces and what it needs,
I'll figure out the right order and run them. I don't care HOW."

## How it works

1. Each module declares itself (PKGBUILD-like):
```bash
   pkgname=libfoo
   produces=(lib/libfoo.a)
   depends=(libbar)
   ...
```

2. foe-hammer builds a dependency graph

3. Each module runs its build however it wants:
   - CMake? ✓
   - Makefile? ✓
   - bash script? ✓
   - Zig build? ✓
   - Cargo? ✓

4. foe-hammer orchestrates the order and cross-compilation targets