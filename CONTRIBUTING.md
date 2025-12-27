# Contributing to Fíth Renderer

Thank you for your interest in contributing to the Fíth template renderer! This document provides guidelines for contributing to ensure code quality and consistency.

## Code Quality Standards

This project maintains high code quality standards through strict adherence to SOLID principles, Go programming best practices, and comprehensive testing.

### SOLID Principles

All code must follow the five SOLID principles:

#### 1. Single Responsibility Principle (SRP)
- Each package has one clear purpose
- Each type has a single reason to change
- Example: The `lexer` package only handles tokenization, not parsing

#### 2. Open/Closed Principle (OCP)
- Software entities should be open for extension but closed for modification
- Use interfaces to allow new behavior without changing existing code
- Example: Custom loaders implement the `Loader` interface

#### 3. Liskov Substitution Principle (LSP)
- Subtypes must be substitutable for their base types
- Any implementation of an interface must be interchangeable
- Example: Any `Loader` implementation can be used without code changes

#### 4. Interface Segregation Principle (ISP)
- Clients should not depend on methods they don't use
- Keep interfaces small and focused
- Example: `Loader` interface has only `Load()` and `Exists()`

#### 5. Dependency Inversion Principle (DIP)
- High-level modules should depend on abstractions, not concretions
- Use dependency injection
- Example: Runtime depends on `Loader` interface, not `DirectoryLoader`

### Go Programming Standards

#### Official Standards
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- All code must pass `gofmt` and `goimports`
- All code must pass `go vet` without warnings
- All code must pass `staticcheck` without warnings

#### Go Idioms
- **Accept interfaces, return structs**
  ```go
  func New(loader Loader) *Renderer { ... }  // Accept interface
  ```
- **Make zero value useful**
  ```go
  var lexer Lexer  // Should be usable without initialization
  ```
- **Error handling is explicit**
  ```go
  result, err := doSomething()
  if err != nil {
      return fmt.Errorf("doing something: %w", err)
  }
  ```
- **Use short, descriptive names**
- **Prefer composition over inheritance**
- **Keep interfaces small**

### Testing Requirements

#### Coverage
- Minimum 90% code coverage across all packages
- Every exported function must have tests
- Use table-driven tests where appropriate

#### Test Types
1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test full pipeline (lexer → parser → runtime)
3. **Benchmarks**: Measure performance of critical paths
4. **Examples**: GoDoc examples that double as tests

#### Test Organization
```go
func TestFeatureName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    Result
        wantErr bool
    }{
        {name: "valid input", input: "test", want: expected},
        {name: "invalid input", input: "", wantErr: true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("unexpected error: %v", err)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Documentation Requirements

#### GoDoc Comments
Every exported symbol must have a GoDoc comment:

```go
// Lexer tokenizes template strings into discrete tokens.
// It tracks line and column positions for error reporting.
type Lexer struct {
    // ...
}

// New creates a new Lexer for the given input string.
// The zero value is not usable; always use New to create instances.
func New(input string) *Lexer {
    // ...
}

// NextToken returns the next token from the input.
// Returns TokenEOF when the end of input is reached.
// Returns an error if invalid syntax is encountered.
func (l *Lexer) NextToken() (Token, error) {
    // ...
}
```

#### Package Documentation
Every package must have a package comment:

```go
// Package lexer provides tokenization for Fíth templates.
//
// The lexer scans template input and produces a stream of tokens
// that can be consumed by the parser. It handles both text content
// and template expressions delimited by {{ and }}.
//
// Example:
//
//	l := lexer.New("Hello {{.Name}}")
//	for {
//	    tok, err := l.NextToken()
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    if tok.Type == lexer.TokenEOF {
//	        break
//	    }
//	    fmt.Printf("%v\n", tok)
//	}
package lexer
```

## Development Workflow

### Before Starting
1. Review relevant specs in `openspec/changes/implement-fith-renderer/specs/`
2. Check tasks.md for your assigned task
3. Read related code to understand context

### Making Changes
1. Create a feature branch: `git checkout -b feature/your-feature`
2. Write tests first (TDD approach)
3. Implement the feature
4. Ensure tests pass: `go test ./...`
5. Run linters: `make lint` or `golangci-lint run`
6. Format code: `gofmt -w .` and `goimports -w .`
7. Update documentation if needed
8. Commit with clear message

### Code Review Checklist

Before submitting a PR, verify:

- [ ] Follows SOLID principles
- [ ] Adheres to Go idioms and conventions
- [ ] `gofmt` and `goimports` applied
- [ ] `go vet` passes
- [ ] `staticcheck` or `golangci-lint` passes
- [ ] Tests added or updated
- [ ] Test coverage >90% for changed code
- [ ] Documentation updated (GoDoc, README, etc.)
- [ ] No unnecessary allocations in hot paths
- [ ] Error handling is explicit and clear
- [ ] Interfaces are minimal and focused
- [ ] Dependencies injected (not hardcoded)

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./benchmarks

# Run specific package tests
go test ./lexer/...
```

### Code Style

#### Naming Conventions
- **Packages**: lowercase, single word (`lexer`, not `lexer_utils`)
- **Exported types**: PascalCase (`Lexer`, `Token`)
- **Unexported types**: camelCase (`tokenBuffer`)
- **Interfaces**: typically noun or adjective (`Loader`, `Renderable`)
- **Functions**: camelCase or PascalCase depending on export

#### Error Handling
```go
// Sentinel errors
var ErrNotFound = errors.New("template not found")

// Custom error types
type SyntaxError struct {
    Line, Column int
    Message      string
}

func (e *SyntaxError) Error() string {
    return fmt.Sprintf("syntax error at %d:%d: %s", e.Line, e.Column, e.Message)
}

// Wrapping errors
if err != nil {
    return fmt.Errorf("parsing template: %w", err)
}
```

#### Comments
- Comment code that needs clarification, not obvious operations
- Explain "why" not "what"
- Keep comments up-to-date with code

## Project Structure

```
toutago-fith-renderer/
├── lexer/          # Tokenization
├── parser/         # AST generation
├── compiler/       # Template compilation
├── runtime/        # Execution engine
├── loader/         # Template loading
├── builtins/       # Built-in functions
├── examples/       # Usage examples
├── docs/           # Documentation
└── benchmarks/     # Performance tests
```

## Questions or Issues?

- Open an issue on GitHub
- Refer to the design document: `openspec/changes/implement-fith-renderer/design.md`
- Review the specification: `openspec/changes/implement-fith-renderer/specs/`

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
