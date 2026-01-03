# orchestrator

Orchestrates module builds without knowing how to compile.

**foe-hammer does NOT compile anything.** It executes hooks defined in PKGBUILDs.

---

## Responsibilities

| Component | Knows | Doesn't know |
|-----------|-------|--------------|
| ModuleLoader | Parse PKGBUILD | What to do with it |
| ModuleGraph | Order, cycles, dependencies | How to build |
| Orchestrator | When to build, in what order | What's inside hooks |
| ContextProvider | Host, target, paths | What module does with it |
| HookRunner | Execute a hook with env | What's inside |

---

## Flow
```
┌─────────────────────────────────────────────────────────┐
│                        LOAD                             │
│                                                         │
│   Load(rootDir)                                         │
│       ├── loader.LoadAll()     → []*Module              │
│       ├── graph.Add(modules)                            │
│       ├── graph.Validate()     → deps exist?            │
│       └── graph.TopoSort()     → order + cycle detect   │
│                                                         │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                     CONFIGURE                           │
│                                                         │
│   SetOutput(outDir)            → where to write files   │
│                                                         │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                        PLAN                             │
│                                                         │
│   Plan(target)                                          │
│       └── for each module:                              │
│               produces = runner.Produces(module, env)   │
│               module.Produces = produces                │
│                                                         │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                        BUILD                            │
│                                                         │
│   BuildAll(target)                                      │
│       └── for name := range Order()                     │
│               Build(name, target)                       │
│                                                         │
│   BuildFrom(name, target)                               │
│       └── for name := range Descendants(name)           │
│               Build(name, target)                       │
│                                                         │
│   Build(name, target)                                   │
│       ├── CanBuild(name)       → check tools            │
│       ├── BuildEnv(...)        → prepare env            │
│       └── runner.Run(m, env)   → execute hook           │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## Usage
```go
host := domain.Host{OS: runtime.GOOS, Arch: runtime.GOARCH}
target := domain.Target{OS: "linux", Arch: "amd64"}

orchestrator := app.NewOrchestrator(
    moduleloader.NewBashLoader(),
    context.NewEnvProvider(),
    hookrunner.NewBashHookRunner(),
    host,
    toolchecker.NewWhichChecker(),
)

// 1. Load
err := orchestrator.Load("./myproject")

// 2. Configure
orchestrator.SetOutput("./build")

// 3. Plan
err = orchestrator.Plan(target)

// 4. Build
err = orchestrator.BuildAll(target)

// Or rebuild from a specific module
err = orchestrator.BuildFrom("libb", target)
```

