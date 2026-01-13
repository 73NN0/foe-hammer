# Notes

## About Architecture

Threedotslabs describes their services with ports and adapters folders.
"Ports" here doesn't mean interfaces like in the pure hexagonal pattern, but inbound adapters (inbound in the sense that the query is coming into the service).
"Adapters" here means outbound adapters (in the sense that the query exits the service to target external services like databases, etc.).

I know they have factories in their code, but I don't think I need them right now.

My mindset is relatively simple:
- First, there is data on one side and code that processes the data on the other side. We need to abstract the data from the code.
- Data can be modeled using the OOP pattern (classic) or a data-driven approach.

For example, let's look at `module.go`, which describes a module inside foe-hammer:
```go
type Module struct {
    Name        string
    DirPath     string   // Directory where the module lives
    Path        string   // Full path (absolute)
    Description string
    Produces    []string // Relative paths of build artifacts
    Depends     []string // Dependency modules
    MakeDepends []string // External dependencies (SDL2, etc.)
    Sources     []string // Source files
}
```

It's simple and OOP. For POC considerations, every property is public (why not? It's O(1) access in both time and space), and the type `Module` describes an "object" module.

But we can consider modules from another architectural point of view:
```go
type Modules struct {
    Id          []int8
    DirPaths    []string
    Paths       []string
    Descriptions []string
    Produces    [][]string
    Depends     [][]string
    MakeDepends [][]string
    Sources     [][]string
}
```

But there are some catches here we need to address.

First, a string is flexible in size by nature. Second, `[][]string` is very expensive in memory for no good reason.

Unlike C, `[][]` is not syntactic sugar but two nested slices.
One slice is 24 bytes:
```go
slice = {
    ptr *data  // 8 bytes
    len int    // 8 bytes
    cap int    // 8 bytes
}
```

So for the structure described above, there is approximately 14.4 KB of overhead:
```
Outer slice:     1 × 24 bytes
Inner slices:  100 × 24 bytes = 2400 bytes (headers only!)
String slices: 500 × 24 bytes = 12000 bytes
```

But if we think about it another way—`[]string` + offsets—we have only ~856 bytes of total overhead, and for multiple modules!

---

## Flat Arrays: Theory and Formulas

There are two cases to consider when flattening multidimensional data into a single array.

### Case 1: Regular Matrix (fixed width)

When all rows have the same size (e.g., a framebuffer, a grid, an image), we use arithmetic formulas.

**Storage (2D → 1D):**
```
index = (row × width) + col
```

**Access (1D → 2D):**
```
row = index / width
col = index % width
```

**Example:** A 4×3 image (4 columns, 3 rows)
```
2D view:                 1D storage:
[0,0] [0,1] [0,2] [0,3]  → [0] [1] [2]  [3]  [4]  [5]  [6]  [7]  [8]  [9]  [10] [11]
[1,0] [1,1] [1,2] [1,3]
[2,0] [2,1] [2,2] [2,3]

Access [1, 2]: index = 1 × 4 + 2 = 6 ✓
Reverse index 6: row = 6 / 4 = 1, col = 6 % 4 = 2 ✓
```

**Use cases:** framebuffers, textures, game grids, any fixed-size 2D data.

---

### Case 2: Irregular Data (variable width)

When rows have different sizes (e.g., sources per module), there is no fixed width. We use **offsets** instead.

**Storage:**
```
Data:    ["main.c", "util.c", "app.c", "lib.c", "api.c", "core.c"]
Offsets: [0,        2,        3,       6]
          ^mod0     ^mod1     ^mod2    ^sentinel
```

**Access (module → data):**
```go
start := Offsets[i]
end := Offsets[i + 1]
items := Data[start:end]
```

**Reverse access (flat index → module):**
```go
// Find the largest offset ≤ index
// Linear: O(n), Binary search: O(log n)
```

**The structure becomes:**
```go
type Modules struct {
    Id                 []int8
    DirPaths           []string
    Paths              []string
    Descriptions       []string
    Produces           []string
    ProducesOffsets    []int
    Depends            []string
    DependsOffsets     []int
    MakeDepends        []string
    MakeDependsOffsets []int
    Sources            []string
    SourcesOffsets     []int
}
```

If you find this a little dirty, you can always use the composition principle.

---

## Operations on Flat Arrays

### Access

**Regular matrix:**
```go
func Get(data []T, row, col, width int) T {
    return data[row * width + col]
}

func Set(data []T, row, col, width int, value T) {
    data[row * width + col] = value
}
```

**Irregular data:**
```go
func GetItems(data []T, offsets []int, index int) []T {
    return data[offsets[index]:offsets[index + 1]]
}

func FindOwner(offsets []int, flatIndex int) int {
    // Binary search: find largest i where offsets[i] <= flatIndex
    lo, hi := 0, len(offsets) - 1
    for lo < hi {
        mid := (lo + hi + 1) / 2
        if offsets[mid] <= flatIndex {
            lo = mid
        } else {
            hi = mid - 1
        }
    }
    return lo
}
```

---

### Deletion

Immediate deletion in a flat array is expensive (requires shifting all subsequent elements). Instead, we use **tombstones**: mark as deleted, compact later.

**Strategy: Mark and Sweep**
```go
type Modules struct {
    Sources        []string
    SourcesOffsets []int
    Deleted        []bool   // Tombstone flags
}

// Step 1: Mark as deleted (O(1))
func MarkDeleted(m *Modules, moduleIndex int) {
    m.Deleted[moduleIndex] = true
}

// Step 2: Compact in batch (O(n), done periodically)
func Compact(m *Modules) {
    newSources := make([]string, 0, len(m.Sources))
    newOffsets := make([]int, 0, len(m.SourcesOffsets))
    newDeleted := make([]bool, 0, len(m.Deleted))
    
    newOffsets = append(newOffsets, 0)
    
    for i := 0; i < len(m.Deleted); i++ {
        if m.Deleted[i] {
            continue // Skip deleted modules
        }
        
        start := m.SourcesOffsets[i]
        end := m.SourcesOffsets[i + 1]
        newSources = append(newSources, m.Sources[start:end]...)
        newOffsets = append(newOffsets, len(newSources))
        newDeleted = append(newDeleted, false)
    }
    
    m.Sources = newSources
    m.SourcesOffsets = newOffsets
    m.Deleted = newDeleted
}
```

**Why this approach?**
- Single deletions: O(1)
- Batch compact: O(n), but done rarely
- Amortized cost is low
- Memory stays contiguous after compaction

---

### Insertion

Insertion at the end is cheap. Insertion in the middle is expensive.

**Append (O(1) amortized):**
```go
func AppendModule(m *Modules, sources []string) {
    m.Sources = append(m.Sources, sources...)
    m.SourcesOffsets = append(m.SourcesOffsets, len(m.Sources))
    m.Deleted = append(m.Deleted, false)
}
```

**Insert in middle (expensive, avoid if possible):**
- Requires shifting all data after insertion point
- Requires updating all offsets after insertion point
- Consider: do you really need ordered insertion?

**Alternative: Insert at end + sort by key if order matters**

---

### Batch Operations

Batch operations amortize the cost of expensive operations.

**Batch insert:**
```go
func BatchAppend(m *Modules, allSources [][]string) {
    for _, sources := range allSources {
        m.Sources = append(m.Sources, sources...)
        m.SourcesOffsets = append(m.SourcesOffsets, len(m.Sources))
        m.Deleted = append(m.Deleted, false)
    }
}
```

**Batch delete + compact:**
```go
func BatchDelete(m *Modules, indices []int) {
    for _, i := range indices {
        m.Deleted[i] = true
    }
    Compact(m) // Single compaction pass
}
```

**When to compact?**
- After N deletions
- When deleted ratio exceeds threshold (e.g., 25%)
- On explicit user request
- Before serialization

---

### Iteration

**Iterate all (skip deleted):**
```go
func ForEach(m *Modules, fn func(index int, sources []string)) {
    for i := 0; i < len(m.Deleted); i++ {
        if m.Deleted[i] {
            continue
        }
        start := m.SourcesOffsets[i]
        end := m.SourcesOffsets[i + 1]
        fn(i, m.Sources[start:end])
    }
}
```

**Iterate single field (cache-friendly):**
```go
// Iterating only Paths is a single contiguous memory scan
for i, path := range m.Paths {
    if m.Deleted[i] {
        continue
    }
    process(path)
}
```

This is where SOA shines: iterating one field doesn't load unrelated data into cache.

---

## Summary

| Operation | Regular Matrix | Irregular (Offsets) |
|-----------|---------------|---------------------|
| Access | `row * width + col` | `Offsets[i]` to `Offsets[i+1]` |
| Reverse | `index / width`, `index % width` | Binary search on offsets |
| Delete | Tombstone or shift | Tombstone + batch compact |
| Insert end | O(1) | O(1) amortized |
| Insert middle | O(n) | O(n), avoid if possible |
| Iteration | Simple loop | Loop + skip deleted |

---

## References

- Casey Muratori - Data-Oriented Design
- Mike Acton - "Data-Oriented Design and C++" (CppCon 2014)
- Pikuma - Game engine courses (matrix formulas)
- Threedotslabs - DDD/Clean Architecture in Go

---

## Architecture DDD Threedotslabs

- Tag 2.3 (DDD Lite + Repository Pattern)
    - https://threedots.tech/post/ddd-lite-in-go-introduction/
    - https://threedots.tech/post/repository-pattern-in-go/
    - https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/tree/v2.3
- Tag 2.4 (Clean Architecture)
    - https://threedots.tech/post/introducing-clean-architecture/
    - https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/tree/v2.4
- CQRS (commit 8d9274811559399461aa9f6bf3829316b8ddfb63)
    - https://threedots.tech/post/basic-cqrs-in-go/
    - https://threedots.tech/post/ddd-cqrs-clean-architecture-combined/
    - https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/tree/v2.5
- Flaws and Anti-patterns
    - https://threedots.tech/post/common-anti-patterns-in-go-web-applications/
    - https://threedots.tech/post/microservices-test-architecture/

## Queue in go


// src
// https://stackoverflow.com/questions/23531891/how-do-i-succinctly-remove-the-first-element-from-a-slice-in-go?utm_source=chatgpt.com
// but need to search before implementing this kind of queue

// q := make([]string, 0, len(g.modules))
// head := 0

// push := func(x string) { q = append(q, x) }
// pop := func() string {
//     x := q[head]
//     head++
//     // option anti “memory retention” si besoin
//     if head > 1024 && head*2 >= len(q) {
//         q = append([]string(nil), q[head:]...)
//         head = 0
//     }
//     return x
// }
// empty := func() bool { return head >= len(q) }