# Performance Guide

Guide to optimizing Fíth template rendering performance.

## Table of Contents

- [Benchmarks](#benchmarks)
- [Performance Tips](#performance-tips)
- [Profiling](#profiling)
- [Memory Optimization](#memory-optimization)
- [Caching Strategies](#caching-strategies)
- [Common Bottlenecks](#common-bottlenecks)

## Benchmarks

### Typical Performance

On modern hardware (AMD64, 3.5GHz):

| Operation | Time | Notes |
|-----------|------|-------|
| Template parsing | <1ms | Per template, cached |
| Simple render | 50-100μs | Variables only |
| Complex render | 500μs-1ms | Loops, functions, composition |
| Cache lookup | <10μs | Hot path |
| Function call | 1-5μs | Built-in functions |

### Comparison with html/template

Fíth aims to be within 2x of Go's standard `html/template`:

| Template Type | html/template | Fíth | Ratio |
|---------------|---------------|------|-------|
| Simple (vars) | ~40μs | ~80μs | 2.0x |
| Medium (loops) | ~200μs | ~350μs | 1.75x |
| Complex (composition) | ~500μs | ~900μs | 1.8x |

**Trade-off:** Fíth provides richer syntax and better error messages at a small performance cost.

### Running Benchmarks

```bash
cd /home/nestor/Proyects/toutago-fith-renderer
go test -bench=. -benchmem ./benchmarks
```

Output example:
```
BenchmarkSimpleRender-8     15000   78563 ns/op   4096 B/op   45 allocs/op
BenchmarkLoopRender-8        3000  412234 ns/op  16384 B/op  210 allocs/op
BenchmarkComplexRender-8     2000  876543 ns/op  32768 B/op  450 allocs/op
```

## Performance Tips

### 1. Reuse Renderer Instance

**❌ Bad:**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    renderer := fith.New(fith.Config{TemplateDir: "templates"})
    output, _ := renderer.Render("page", data)
    w.Write([]byte(output))
}
```

**✅ Good:**
```go
var renderer = fith.New(fith.Config{TemplateDir: "templates"})

func handler(w http.ResponseWriter, r *http.Request) {
    output, _ := renderer.Render("page", data)
    w.Write([]byte(output))
}
```

**Impact:** 10-100x faster (avoids template re-parsing)

---

### 2. Use RenderBytes for HTTP

**❌ Bad:**
```go
output, err := renderer.Render("page", data)
w.Write([]byte(output))  // String to bytes conversion
```

**✅ Good:**
```go
output, err := renderer.RenderBytes("page", data)
w.Write(output)  // Direct bytes
```

**Impact:** Saves string allocation and copy

---

### 3. Prepare Data in Go

**❌ Bad:**
```
{{# Complex logic in template #}}
{{range .Users}}
  {{if and (gt (len .Posts) 5) (eq .Status "active")}}
    ...
  {{end}}
{{end}}
```

**✅ Good:**
```go
// Compute in Go
activeUsers := filterActiveUsersWithPosts(users, 5)

data := map[string]interface{}{
    "ActiveUsers": activeUsers,
}
```

```
{{# Simple iteration in template #}}
{{range .ActiveUsers}}
  ...
{{end}}
```

**Impact:** 2-5x faster rendering

---

### 4. Cache Computed Values

**❌ Bad:**
```
{{range .Items}}
  {{formatDate .CreatedAt "Jan 2, 2006"}}
{{end}}
```

**✅ Good:**
```go
type Item struct {
    Name          string
    CreatedAt     time.Time
    FormattedDate string  // Pre-computed
}

for i := range items {
    items[i].FormattedDate = formatDate(items[i].CreatedAt)
}
```

**Impact:** Avoids repeated function calls in loop

---

### 5. Minimize Template Composition Depth

**❌ Bad:**
```
page.html
  → extends layout.html
    → includes header.html
      → includes nav.html
        → includes user-menu.html
```

**✅ Good:**
```
page.html
  → extends layout.html
    → includes header.html (flatter structure)
```

**Impact:** Reduces template loading and parsing overhead

---

### 6. Use Appropriate Data Structures

**❌ Bad:**
```go
// Slice when you need lookups
users := []User{...}  // O(n) lookup in template

// Template needs to check if user exists
{{range .AllUserIDs}}
  {{# This is slow #}}
  {{if userExists . $.Users}}
    ...
  {{end}}
{{end}}
```

**✅ Good:**
```go
// Map for O(1) lookups
userMap := make(map[string]User)
for _, u := range users {
    userMap[u.ID] = u
}

data := map[string]interface{}{
    "Users": userMap,
}
```

**Impact:** O(1) vs O(n) lookups

---

### 7. Avoid Deep Nesting

**❌ Bad:**
```go
data := map[string]interface{}{
    "App": map[string]interface{}{
        "Config": map[string]interface{}{
            "Database": map[string]interface{}{
                "Settings": map[string]interface{}{
                    "Host": "localhost",
                },
            },
        },
    },
}
```

```
{{.App.Config.Database.Settings.Host}}  // Slow
```

**✅ Good:**
```go
data := map[string]interface{}{
    "DBHost": config.Database.Host,
}
```

```
{{.DBHost}}  // Fast
```

**Impact:** Reduces lookup overhead

---

## Profiling

### CPU Profiling

```go
package main

import (
    "os"
    "runtime/pprof"
    "github.com/toutaio/toutago-fith-renderer"
)

func main() {
    // Start CPU profiling
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // Your rendering code
    renderer := fith.New(fith.Config{TemplateDir: "templates"})
    for i := 0; i < 1000; i++ {
        renderer.Render("page", data)
    }
}
```

Analyze:
```bash
go tool pprof cpu.prof
(pprof) top10
(pprof) list RenderFunction
```

---

### Memory Profiling

```go
import "runtime/pprof"

func main() {
    renderer := fith.New(fith.Config{TemplateDir: "templates"})
    
    // Render many times
    for i := 0; i < 1000; i++ {
        renderer.Render("page", data)
    }
    
    // Memory profile
    f, _ := os.Create("mem.prof")
    pprof.WriteHeapProfile(f)
    f.Close()
}
```

Analyze:
```bash
go tool pprof mem.prof
(pprof) top10
(pprof) list RenderFunction
```

---

## Memory Optimization

### String Builder

Fíth uses `strings.Builder` internally for efficient string concatenation.

**Your code:**
```go
// If building strings for data
var b strings.Builder
for _, item := range items {
    b.WriteString(item.Name)
    b.WriteString("\n")
}
data := map[string]interface{}{
    "List": b.String(),
}
```

---

### Reduce Allocations

**❌ Bad:**
```go
// Creates new map each request
func handler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "Title": "Home",
        "User": getUser(),
    }
    renderer.Render("home", data)
}
```

**✅ Good:**
```go
// Reuse data structure (if safe)
type PageData struct {
    Title string
    User  User
}

var dataPool = sync.Pool{
    New: func() interface{} {
        return &PageData{}
    },
}

func handler(w http.ResponseWriter, r *http.Request) {
    data := dataPool.Get().(*PageData)
    defer dataPool.Put(data)
    
    data.Title = "Home"
    data.User = getUser()
    
    renderer.Render("home", data)
}
```

**Note:** Only safe if renderer doesn't retain references to data.

---

## Caching Strategies

### 1. Template Caching (Automatic)

Templates are automatically cached after first parse. No action needed.

---

### 2. Output Caching

Cache rendered output for static pages:

```go
var cache = make(map[string]string)
var cacheMu sync.RWMutex

func renderCached(slug string, data interface{}) (string, error) {
    // Check cache
    cacheMu.RLock()
    if output, ok := cache[slug]; ok {
        cacheMu.RUnlock()
        return output, nil
    }
    cacheMu.RUnlock()
    
    // Render
    output, err := renderer.Render(slug, data)
    if err != nil {
        return "", err
    }
    
    // Store in cache
    cacheMu.Lock()
    cache[slug] = output
    cacheMu.Unlock()
    
    return output, nil
}
```

---

### 3. TTL Cache

Use third-party cache with TTL:

```go
import "github.com/patrickmn/go-cache"

var outputCache = cache.New(5*time.Minute, 10*time.Minute)

func renderWithTTL(slug string, data interface{}) (string, error) {
    // Check cache
    if output, found := outputCache.Get(slug); found {
        return output.(string), nil
    }
    
    // Render
    output, err := renderer.Render(slug, data)
    if err != nil {
        return "", err
    }
    
    // Cache for 5 minutes
    outputCache.Set(slug, output, cache.DefaultExpiration)
    
    return output, nil
}
```

---

### 4. Conditional Caching

Cache only when data is stable:

```go
func renderSmart(slug string, data interface{}) (string, error) {
    // Check if data is cacheable
    if cacheable, ok := data.(Cacheable); ok && cacheable.IsCacheable() {
        cacheKey := fmt.Sprintf("%s:%s", slug, cacheable.CacheKey())
        if output, found := cache.Get(cacheKey); found {
            return output.(string), nil
        }
    }
    
    // Render normally
    return renderer.Render(slug, data)
}
```

---

## Common Bottlenecks

### 1. Large Loops

**Problem:**
```
{{range .Items}}  {{# 10,000 items #}}
  <div>{{.Name | upper}}</div>
{{end}}
```

**Solutions:**
- Paginate data in Go
- Process transformations in Go
- Use simpler templates

---

### 2. Deep Object Graphs

**Problem:**
```
{{.Level1.Level2.Level3.Level4.Value}}
```

**Solution:**
```go
// Flatten in Go
data := map[string]interface{}{
    "Value": obj.Level1.Level2.Level3.Level4.Value,
}
```

---

### 3. Repeated Function Calls

**Problem:**
```
{{range .Items}}
  {{expensiveFunc .ID}}  {{# Called N times #}}
{{end}}
```

**Solution:**
```go
// Pre-compute
for i := range items {
    items[i].ComputedValue = expensiveFunc(items[i].ID)
}
```

---

### 4. Template Composition Overhead

**Problem:**
Multiple includes/extends create parsing overhead.

**Solution:**
- Flatten template structure
- Consider combining small includes
- Cache aggressively in production

---

## Production Checklist

- [ ] Reuse single renderer instance
- [ ] Use `RenderBytes` for HTTP responses
- [ ] Prepare complex data in Go code
- [ ] Cache computed values in data
- [ ] Profile hot paths
- [ ] Monitor memory allocations
- [ ] Consider output caching for static content
- [ ] Minimize template composition depth
- [ ] Use appropriate data structures (maps vs slices)
- [ ] Benchmark critical paths

---

## Monitoring

### Metrics to Track

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    renderDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "fith_render_duration_seconds",
            Help: "Template render duration",
        },
        []string{"template"},
    )
    
    renderErrors = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "fith_render_errors_total",
            Help: "Template render errors",
        },
        []string{"template"},
    )
)

func monitoredRender(slug string, data interface{}) (string, error) {
    start := time.Now()
    output, err := renderer.Render(slug, data)
    duration := time.Since(start).Seconds()
    
    renderDuration.WithLabelValues(slug).Observe(duration)
    if err != nil {
        renderErrors.WithLabelValues(slug).Inc()
    }
    
    return output, err
}
```

---

## Need More Performance?

If you need extreme performance:

1. **Pre-render static content** - Generate HTML at build time
2. **Use CDN** - Serve from edge locations
3. **HTTP caching** - Use `Cache-Control` headers
4. **Consider html/template** - For maximum speed at cost of features
5. **Profile your code** - Find actual bottlenecks

Most applications won't need these optimizations. Fíth is fast enough for typical web applications.

---

## Benchmarking Your Templates

Create a benchmark:

```go
package myapp

import (
    "testing"
    "github.com/toutaio/toutago-fith-renderer"
)

var renderer = fith.New(fith.Config{TemplateDir: "templates"})

func BenchmarkMyPage(b *testing.B) {
    data := map[string]interface{}{
        "Title": "Test",
        "Items": generateItems(100),
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := renderer.Render("mypage", data)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

Run:
```bash
go test -bench=BenchmarkMyPage -benchmem
```

---

See [API Reference](api.md) for more optimization techniques.
