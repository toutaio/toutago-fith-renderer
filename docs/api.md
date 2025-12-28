# Fíth API Reference

Complete API documentation for using Fíth in your Go applications.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Rendering Templates](#rendering-templates)
- [Custom Functions](#custom-functions)
- [Template Loading](#template-loading)
- [Error Handling](#error-handling)
- [Performance](#performance)

## Installation

```bash
go get github.com/toutaio/toutago-fith-renderer
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/toutaio/toutago-fith-renderer"
)

func main() {
    // Create a new Fíth renderer
    renderer := fith.New(fith.Config{
        TemplateDir: "templates",
    })
    
    // Prepare your data
    data := map[string]interface{}{
        "Title": "Welcome",
        "User": map[string]interface{}{
            "Name": "Alice",
            "Email": "alice@example.com",
        },
    }
    
    // Render a template
    output, err := renderer.Render("home", data)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(output)
}
```

## Configuration

### Config Structure

```go
type Config struct {
    // TemplateDir is the root directory for templates
    TemplateDir string
    
    // TemplateExt is the file extension for templates (default: ".html")
    TemplateExt string
    
    // Loader is a custom template loader (optional)
    Loader Loader
    
    // Functions is a map of custom functions to register
    Functions map[string]runtime.Function
    
    // StrictMode enables strict variable checking (future feature)
    StrictMode bool
}
```

### Basic Configuration

```go
renderer := fith.New(fith.Config{
    TemplateDir: "templates",
    TemplateExt: ".html",
})
```

### With Custom Functions

```go
renderer := fith.New(fith.Config{
    TemplateDir: "templates",
    Functions: map[string]runtime.Function{
        "greet": func(args ...interface{}) (interface{}, error) {
            if len(args) != 1 {
                return nil, fmt.Errorf("greet expects 1 argument")
            }
            name := args[0].(string)
            return fmt.Sprintf("Hello, %s!", name), nil
        },
    },
})
```

### With Custom Loader

```go
loader := loader.NewDirectoryLoader("./templates", ".html")
renderer := fith.New(fith.Config{
    Loader: loader,
})
```

## Rendering Templates

### Render Method

```go
func (f *Fith) Render(slug string, data interface{}) (string, error)
```

Renders a template with the given data.

**Parameters:**
- `slug`: Template identifier (e.g., "home", "user/profile")
- `data`: Data to pass to the template (map, struct, or any value)

**Returns:**
- `string`: Rendered output
- `error`: Error if rendering fails

### Example: Render with Map

```go
data := map[string]interface{}{
    "Title": "My Page",
    "Items": []string{"one", "two", "three"},
}

output, err := renderer.Render("page", data)
```

### Example: Render with Struct

```go
type User struct {
    Name  string
    Email string
    Admin bool
}

user := User{
    Name:  "Bob",
    Email: "bob@example.com",
    Admin: true,
}

output, err := renderer.Render("user-profile", user)
```

### RenderBytes Method

```go
func (f *Fith) RenderBytes(slug string, data interface{}) ([]byte, error)
```

Like Render, but returns bytes instead of string (more efficient for HTTP responses).

```go
data := map[string]interface{}{"Title": "Home"}
output, err := renderer.RenderBytes("home", data)
if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
w.Write(output)
```

## Custom Functions

### Function Signature

```go
type Function func(args ...interface{}) (interface{}, error)
```

### Register Functions at Creation

```go
renderer := fith.New(fith.Config{
    TemplateDir: "templates",
    Functions: map[string]runtime.Function{
        "double": func(args ...interface{}) (interface{}, error) {
            if len(args) != 1 {
                return nil, fmt.Errorf("double expects 1 argument")
            }
            n, ok := args[0].(int)
            if !ok {
                return nil, fmt.Errorf("double expects an integer")
            }
            return n * 2, nil
        },
    },
})
```

### Register Functions After Creation

```go
renderer := fith.New(fith.Config{TemplateDir: "templates"})

renderer.RegisterFunction("myFunc", func(args ...interface{}) (interface{}, error) {
    // Your function logic
    return "result", nil
})
```

### Type-Safe Function Example

```go
func createGreetFunction() runtime.Function {
    return func(args ...interface{}) (interface{}, error) {
        // Validate argument count
        if len(args) != 1 {
            return nil, fmt.Errorf("greet expects 1 argument, got %d", len(args))
        }
        
        // Type check
        name, ok := args[0].(string)
        if !ok {
            return nil, fmt.Errorf("greet expects a string argument")
        }
        
        // Return result
        return fmt.Sprintf("Hello, %s!", name), nil
    }
}

renderer.RegisterFunction("greet", createGreetFunction())
```

## Template Loading

### Directory Loader

Load templates from the filesystem:

```go
loader := loader.NewDirectoryLoader("./templates", ".html")
renderer := fith.New(fith.Config{Loader: loader})
```

**Slug Resolution:**
- `"home"` → `./templates/home.html`
- `"user/profile"` → `./templates/user/profile.html`
- `"admin/users/list"` → `./templates/admin/users/list.html`

### Embed Loader

Embed templates in your binary using Go 1.16+ embed:

```go
//go:embed templates/*
var templateFS embed.FS

func main() {
    loader := loader.NewEmbedLoader(templateFS, "templates", ".html")
    renderer := fith.New(fith.Config{Loader: loader})
}
```

### Custom Loader

Implement the Loader interface:

```go
type Loader interface {
    Load(slug string) (string, error)
}
```

Example custom loader:

```go
type DatabaseLoader struct {
    db *sql.DB
}

func (l *DatabaseLoader) Load(slug string) (string, error) {
    var content string
    err := l.db.QueryRow("SELECT content FROM templates WHERE slug = ?", slug).Scan(&content)
    if err != nil {
        return "", fmt.Errorf("template not found: %s", slug)
    }
    return content, nil
}

// Use it
loader := &DatabaseLoader{db: db}
renderer := fith.New(fith.Config{Loader: loader})
```

### Caching

Templates are automatically cached after first load. To disable caching (e.g., in development):

```go
// Note: Cache control API is planned for future release
renderer := fith.New(fith.Config{
    TemplateDir: "templates",
    // NoCache: true, // Future feature
})
```

## Error Handling

### Error Types

```go
// TemplateError represents an error during template processing
type TemplateError struct {
    Template string  // Template name/slug
    Line     int     // Line number
    Column   int     // Column number
    Message  string  // Error message
    Err      error   // Underlying error
}

func (e *TemplateError) Error() string
```

### Handling Errors

```go
output, err := renderer.Render("home", data)
if err != nil {
    // Type assertion for detailed error info
    if te, ok := err.(*fith.TemplateError); ok {
        log.Printf("Template error in %s at line %d, column %d: %s",
            te.Template, te.Line, te.Column, te.Message)
    } else {
        log.Printf("Error: %v", err)
    }
    return
}
```

### Common Errors

**Template Not Found:**
```go
output, err := renderer.Render("nonexistent", data)
// Error: template not found: nonexistent
```

**Unknown Variable:**
```go
// Template: {{.UnknownField}}
output, err := renderer.Render("home", data)
// Error: unknown variable: .UnknownField at line 5, column 3
```

**Unknown Function:**
```go
// Template: {{unknownFunc .Name}}
output, err := renderer.Render("home", data)
// Error: unknown function: unknownFunc at line 8, column 5
```

## Performance

### Benchmarks

Typical performance (on modern hardware):

- **Template parsing:** <1ms for typical templates
- **Rendering (simple):** ~50-100μs
- **Rendering (complex):** ~500μs-1ms
- **Cache lookup:** <10μs

### Optimization Tips

#### 1. Reuse Renderer Instance

```go
// Good: Create once, reuse many times
var renderer = fith.New(fith.Config{TemplateDir: "templates"})

func handler(w http.ResponseWriter, r *http.Request) {
    output, _ := renderer.Render("page", data)
    w.Write([]byte(output))
}
```

#### 2. Use RenderBytes for HTTP

```go
// More efficient: avoids string conversion
output, err := renderer.RenderBytes("page", data)
w.Write(output)
```

#### 3. Prepare Data in Go

```go
// Good: Complex logic in Go
data := map[string]interface{}{
    "UserCount": len(users),
    "IsEligible": checkEligibility(user),
}

// Bad: Complex logic in template
// {{if and (gt (len .Users) 0) (eq .User.Status "active")}}
```

#### 4. Cache Computed Values

```go
// Cache expensive computations
data := map[string]interface{}{
    "FormattedDate": formatDate(time.Now()),
    "Summary": generateSummary(article),
}
```

### Memory Usage

Typical memory usage:
- **Cached template:** ~5-10KB per template
- **Rendering context:** ~1-2KB per render
- **Output buffer:** Depends on output size

## Thread Safety

The Fíth renderer is **thread-safe** and can be used concurrently:

```go
var renderer = fith.New(fith.Config{TemplateDir: "templates"})

// Safe to call from multiple goroutines
go func() {
    output, _ := renderer.Render("page1", data1)
}()

go func() {
    output, _ := renderer.Render("page2", data2)
}()
```

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    
    "github.com/toutaio/toutago-fith-renderer"
    "github.com/toutaio/toutago-fith-renderer/runtime"
)

var renderer *fith.Fith

func init() {
    renderer = fith.New(fith.Config{
        TemplateDir: "templates",
        TemplateExt: ".html",
        Functions: map[string]runtime.Function{
            "formatTime": func(args ...interface{}) (interface{}, error) {
                if len(args) != 1 {
                    return nil, fmt.Errorf("formatTime expects 1 argument")
                }
                t, ok := args[0].(time.Time)
                if !ok {
                    return nil, fmt.Errorf("formatTime expects a time.Time")
                }
                return t.Format("Jan 2, 2006 at 3:04pm"), nil
            },
        },
    })
}

type User struct {
    Name      string
    Email     string
    JoinedAt  time.Time
    IsAdmin   bool
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    user := User{
        Name:     "Alice Johnson",
        Email:    "alice@example.com",
        JoinedAt: time.Now().AddDate(-1, 0, 0),
        IsAdmin:  true,
    }
    
    data := map[string]interface{}{
        "Title": "Welcome Home",
        "User":  user,
        "Posts": []string{"First Post", "Second Post", "Third Post"},
    }
    
    output, err := renderer.RenderBytes("home", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write(output)
}

func main() {
    http.HandleFunc("/", homeHandler)
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Next Steps

- [Syntax Reference](syntax.md) - Complete template syntax guide
- [Built-in Functions](functions.md) - All available functions
- [Examples](../examples/) - Working code examples
- [Migration Guide](migration.md) - Migrating from other template engines
