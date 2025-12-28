# Basic Fíth Example

This example demonstrates the basic usage of the Fíth template engine.

## Running the Example

```bash
go run main.go
```

## Features Demonstrated

1. **Simple Variables** - Basic variable substitution with `{{.Name}}`
2. **Conditionals** - If/else statements for conditional rendering
3. **Loops** - Iterating over collections with `{{range}}`
4. **Built-in Functions** - Using built-in functions like `upper`, `lower`, `join`
5. **Filter Pipelines** - Chaining functions with the pipe operator `|`
6. **Custom Functions** - Registering and using custom functions
7. **Nested Data** - Accessing nested data structures with dot notation

## Code Overview

The example creates a Fíth engine with default settings:

```go
engine, err := fith.NewWithDefaults()
```

Then demonstrates various features using `RenderString()` to render template strings directly:

```go
output, err := engine.RenderString(template, data)
```

## Next Steps

- See the `layouts` example for template composition
- See the `filters` example for advanced function usage
- See the `includes` example for template reuse
