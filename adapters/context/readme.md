# provider

ContextProvider builds the environment variables for hook execution
a simple map[string]string.

For now it only use shell environment variables.

## Injected into hooks by ContextProvider:

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