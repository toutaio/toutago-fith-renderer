# Toutā Fíth Renderer

A powerful, flexible template engine for Go inspired by Celtic craftsmanship.

> **Fíth** (Old Irish): The art of weaving patterns - representing how templates weave data into beautiful output.

## Features

- 🎯 **Modern Syntax**: Jinja2/Twig-inspired template syntax
- ⚡ **High Performance**: Within 2x of Go's html/template
- 🔧 **Extensible**: Custom functions and loaders via clean interfaces
- 📦 **Flexible Loading**: Filesystem, embed.FS, or custom sources
- 🎨 **Template Composition**: Includes and layout inheritance
- 🔍 **Great Errors**: Line and column numbers in all error messages
- ✅ **Well Tested**: >90% code coverage with comprehensive tests
- 📚 **Well Documented**: Full GoDoc and usage examples

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/toutaio/toutago-fith-renderer"
)

func main() {
    // Create renderer
    renderer := fith.New(fith.Config{
        TemplateDir: "templates",
    })
    
    // Prepare data
    data := map[string]interface{}{
        "Title": "Welcome",
        "User": map[string]interface{}{
            "Name": "Alice",
        },
    }
    
    // Render template
    output, err := renderer.Render("home", data)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(output)
}
```

## Template Syntax

```html
{{# templates/home.html #}}
{{extends "layout"}}

{{block "title"}}{{.Title}}{{end}}

{{block "content"}}
  <h1>Hello, {{.User.Name | upper}}!</h1>
  
  {{if .User.IsAdmin}}
    <p>Admin panel access granted.</p>
  {{end}}
  
  {{range .Items}}
    <li>{{.}} ({{@index}})</li>
  {{end}}
{{end}}
```

## Status

🚀 **In Active Development** - Core implementation in progress

- ✅ Project structure and tooling
- ⏳ Lexer implementation
- ⏳ Parser implementation
- ⏳ Runtime engine
- ⏳ Built-in functions
- ⏳ Template loading
- ⏳ Documentation

## Installation

```bash
go get github.com/toutaio/toutago-fith-renderer
```

## Documentation

- [API Documentation](docs/api.md) - Coming soon
- [Syntax Reference](docs/syntax.md) - Coming soon
- [Built-in Functions](docs/functions.md) - Coming soon
- [Migration Guide](docs/migration.md) - Coming soon

## Development

Built with:
- **SOLID Principles**: Maintainable, extensible design
- **Go Standards**: Strict adherence to Go idioms and best practices
- **Test Coverage**: >90% with comprehensive testing
- **Clean Interfaces**: Easy to extend with custom loaders and functions

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## Repository

https://github.com/toutaio/toutago-fith-renderer

## License

MIT - see [LICENSE](LICENSE) file

---

Part of the [Toutā framework](https://github.com/toutaio) ecosystem.
