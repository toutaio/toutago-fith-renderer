# Fíth Quick Reference

> One-page reference for the Fíth template language

## Basic Syntax

```
{{.Variable}}              Output a variable
{{# Comment #}}            Comment (not rendered)
```

## Variables

```
{{.Name}}                  Simple variable
{{.User.Name}}             Nested field
{{.User.Profile.Email}}    Deep nesting
```

## Control Flow

### If Statements
```
{{if .IsLoggedIn}}
  Welcome back!
{{end}}

{{if .User.IsAdmin}}
  Admin
{{else}}
  User
{{end}}

{{if .Role eq "admin"}}
  Admin
{{elseif .Role eq "mod"}}
  Moderator
{{else}}
  User
{{end}}
```

### Range Loops
```
{{range .Items}}
  {{.}}
{{end}}

{{range .Items}}
  Item {{@index}}: {{.}}
{{end}}

{{range .Items}}
  {{.Name}} - {{.Price}}
{{else}}
  No items found
{{end}}
```

### Loop Variables
```
{{@index}}    Current index (0-based)
{{@first}}    True if first iteration
{{@last}}     True if last iteration
{{@odd}}      True if odd (1st, 3rd, 5th...)
{{@even}}     True if even (2nd, 4th, 6th...)
```

## Functions

### Function Calls
```
{{upper .Name}}
{{truncate .Text 100}}
{{replace .Content "old" "new"}}
```

### Filter Pipeline
```
{{.Name | upper}}
{{.Name | upper | trim}}
{{.Text | truncate 100 | htmlEscape}}
```

## Built-in Functions

### String Functions
```
{{upper "hello"}}              → HELLO
{{lower "WORLD"}}              → world
{{title "hello world"}}        → Hello World
{{trim "  text  "}}           → text
{{truncate "Long text" 5}}     → Long...
{{replace "Hi Bob" "Bob" "Alice"}}  → Hi Alice
```

### Array Functions
```
{{join .Tags ", "}}           → tag1, tag2, tag3
{{len .Items}}                → 5
{{first .Items}}              → (first item)
{{last .Items}}               → (last item)
```

### Logic Functions
```
{{default .Name "Anonymous"}}  → Anonymous (if Name empty)
```

### Encoding Functions
```
{{urlEncode "hello world"}}    → hello+world
{{htmlEscape "<script>"}}      → &lt;script&gt;
```

### Date Functions
```
{{date "Jan 2, 2006" .CreatedAt}}  → Dec 27, 2024
{{date "15:04" .Time}}             → 14:30
```

## Template Composition

### Include
```
{{include "header"}}
{{include "partials/nav"}}
{{include "user-card" .User}}
```

### Extends & Blocks
```
{{# Layout: layouts/base.html #}}
<!DOCTYPE html>
<html>
<head>
  <title>{{block "title"}}Default{{end}}</title>
</head>
<body>
  {{block "content"}}{{end}}
</body>
</html>

{{# Page: home.html #}}
{{extends "layouts/base"}}

{{block "title"}}Home Page{{end}}

{{block "content"}}
  <h1>Welcome!</h1>
{{end}}
```

## Common Patterns

### User Authentication
```
{{if .User}}
  Welcome, {{.User.Name}}!
  {{if .User.IsAdmin}}
    <a href="/admin">Admin</a>
  {{end}}
{{else}}
  <a href="/login">Login</a>
{{end}}
```

### List with Empty State
```
{{range .Items}}
  <li>{{.}}</li>
{{else}}
  <p>No items found</p>
{{end}}
```

### Alternating Row Styles
```
{{range .Items}}
  <tr class="{{if @odd}}odd{{else}}even{{end}}">
    <td>{{.Name}}</td>
  </tr>
{{end}}
```

### First/Last Indicators
```
{{range .Items}}
  {{if @first}}<strong>{{end}}
  {{.}}
  {{if @last}}</strong>{{end}}
{{end}}
```

### Breadcrumbs
```
{{range .Breadcrumbs}}
  {{if not @first}} > {{end}}
  <a href="{{.URL}}">{{.Title}}</a>
{{end}}
```

## Go API

### Create Renderer
```go
import "github.com/toutaio/toutago-fith-renderer"

renderer := fith.New(fith.Config{
    TemplateDir: "templates",
})
```

### Render Template
```go
data := map[string]interface{}{
    "Title": "Home",
    "User": user,
}

output, err := renderer.Render("home", data)
if err != nil {
    log.Fatal(err)
}
```

### Custom Function
```go
renderer := fith.New(fith.Config{
    TemplateDir: "templates",
    Functions: map[string]runtime.Function{
        "greet": func(args ...interface{}) (interface{}, error) {
            name := args[0].(string)
            return fmt.Sprintf("Hello, %s!", name), nil
        },
    },
})
```

### Use with HTTP
```go
func handler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{"Title": "Home"}
    output, err := renderer.RenderBytes("home", data)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.Write(output)
}
```

## Truthiness

These values are **false**:
- `nil`
- `false`
- `0`
- `""`
- Empty slice/array
- Empty map

All others are **true**.

## Tips

1. **Prepare data in Go** - Keep templates simple
2. **Use structs** - Type-safe data structures
3. **Cache computed values** - Avoid expensive functions in loops
4. **Reuse renderer** - Create once, use many times
5. **Use RenderBytes** - More efficient for HTTP

## File Organization

```
templates/
├── layouts/
│   └── base.html
├── pages/
│   ├── home.html
│   └── about.html
└── partials/
    ├── header.html
    └── footer.html
```

## Error Messages

Fíth provides detailed errors:
```
Error rendering template "home" at line 12, column 5:
  unknown variable: .User.Namee
```

## Documentation

- [API Reference](docs/api.md)
- [Syntax Guide](docs/syntax.md)
- [Functions Reference](docs/functions.md)
- [Migration Guide](docs/migration.md)
- [Performance Guide](docs/performance.md)

---

**Fíth** - The art of weaving patterns ✨
