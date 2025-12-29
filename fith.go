// Package fith provides a powerful, flexible template engine for Go.
//
// Fíth (Old Irish: "The art of weaving patterns") is a template engine inspired by
// Jinja2 and Twig, designed for generating HTML, text, and other formats from templates.
//
// Basic Usage:
//
//	import "github.com/toutaio/toutago-fith-renderer"
//
//	// Create a new engine
//	engine := fith.New(fith.Config{
//	    TemplateDir: "templates",
//	})
//
//	// Render a template
//	data := map[string]interface{}{
//	    "Title": "Welcome",
//	    "User":  user,
//	}
//	output, err := engine.Render("home", data)
//
// Template Syntax:
//
//	Variables:       {{.Name}}
//	Conditionals:    {{if .IsActive}}...{{end}}
//	Loops:           {{range .Items}}...{{end}}
//	Functions:       {{upper .Name}}
//	Filters:         {{.Name | upper | trim}}
//	Includes:        {{include "header"}}
//	Layouts:         {{extends "layout"}} {{block "content"}}...{{end}}
package fith

import (
	"fmt"
	"io/fs"
	"sync"

	"github.com/toutaio/toutago-fith-renderer/compiler"
	"github.com/toutaio/toutago-fith-renderer/lexer"
	"github.com/toutaio/toutago-fith-renderer/loader"
	"github.com/toutaio/toutago-fith-renderer/parser"
	"github.com/toutaio/toutago-fith-renderer/runtime"
)

// Engine is the main Fíth template engine.
type Engine struct {
	config    Config
	loader    loader.Loader
	compiler  *compiler.Compiler
	functions *runtime.FunctionRegistry
	mu        sync.RWMutex
}

// New creates a new Fíth template engine with the given configuration.
func New(config Config) (*Engine, error) {
	// Apply defaults
	config.applyDefaults()

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	engine := &Engine{
		config:    config,
		functions: runtime.NewFunctionRegistry(),
	}

	// Initialize loader based on config
	if err := engine.initializeLoader(); err != nil {
		return nil, WrapError(ErrorTypeLoader, "failed to initialize loader", err)
	}

	// Initialize compiler
	engine.compiler = compiler.NewCompiler(engine.loader)

	return engine, nil
}

// NewWithDefaults creates a new Fíth engine with default configuration.
func NewWithDefaults() (*Engine, error) {
	return New(DefaultConfig())
}

// initializeLoader sets up the template loader based on configuration.
func (e *Engine) initializeLoader() error {
	if e.config.TemplateFS != nil {
		// Use embedded filesystem
		e.loader = loader.NewEmbedLoader(e.config.TemplateFS, ".", e.config.Extensions)
	} else {
		// Use directory loader
		e.loader = loader.NewFileSystemLoader(e.config.TemplateDir, e.config.Extensions)
	}
	return nil
}

// Render renders a template with the given data.
//
// The slug parameter identifies the template (e.g., "layouts/main", "partials/header").
// The data parameter provides the context for template variables.
//
// Example:
//
//	data := map[string]interface{}{
//	    "Title": "Home Page",
//	    "User": user,
//	}
//	html, err := engine.Render("home", data)
func (e *Engine) Render(slug string, data interface{}) (string, error) {
	// Compile the template
	compiled, err := e.compile(slug)
	if err != nil {
		return "", WrapError(ErrorTypeCompilation, fmt.Sprintf("failed to compile template '%s'", slug), err)
	}

	// Create runtime context
	ctx := runtime.NewContext(data)
	ctx.Set("@slug", slug)

	// Create runtime and register functions
	rt := runtime.NewRuntime(ctx)
	e.copyFunctionsToRuntime(rt)

	// Execute the template
	output, err := e.execute(rt, compiled.AST)
	if err != nil {
		return "", WrapError(ErrorTypeRuntime, fmt.Sprintf("failed to execute template '%s'", slug), err)
	}

	return output, nil
}

// RenderString renders a template string directly without loading from a file.
//
// Example:
//
//	html, err := engine.RenderString("Hello {{.Name}}!", map[string]interface{}{
//	    "Name": "World",
//	})
func (e *Engine) RenderString(template string, data interface{}) (string, error) {
	// Parse the template
	tmpl, err := e.parseString(template)
	if err != nil {
		return "", WrapError(ErrorTypeTemplate, "failed to parse template string", err)
	}

	// Create runtime context
	ctx := runtime.NewContext(data)

	// Create runtime and register functions
	rt := runtime.NewRuntime(ctx)
	e.copyFunctionsToRuntime(rt)

	// Execute the template
	output, err := e.execute(rt, tmpl)
	if err != nil {
		return "", WrapError(ErrorTypeRuntime, "failed to execute template string", err)
	}

	return output, nil
}

// RegisterFunction registers a custom function that can be used in templates.
//
// Functions can be called in templates like: {{myFunc .Value}} or {{.Value | myFunc}}
//
// Example:
//
//	engine.RegisterFunction("reverse", func(args ...interface{}) (interface{}, error) {
//	    if len(args) != 1 {
//	        return nil, fmt.Errorf("reverse expects 1 argument")
//	    }
//	    s := fmt.Sprint(args[0])
//	    runes := []rune(s)
//	    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
//	        runes[i], runes[j] = runes[j], runes[i]
//	    }
//	    return string(runes), nil
//	})
func (e *Engine) RegisterFunction(name string, fn runtime.Function) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.functions.Register(name, fn)
}

// ClearCache clears all compiled template caches.
func (e *Engine) ClearCache() {
	e.compiler.ClearCache()
}

// Exists checks if a template exists without loading it.
func (e *Engine) Exists(slug string) bool {
	return e.loader.Exists(slug)
}

// compile compiles a template using the compiler with caching.
func (e *Engine) compile(slug string) (*compiler.CompiledTemplate, error) {
	if !e.config.CacheEnabled {
		// Load and compile without caching
		tmpl, err := e.loader.Load(slug)
		if err != nil {
			return nil, err
		}
		return e.compiler.CompileWithoutCache(tmpl)
	}

	// Use compiler's caching
	return e.compiler.Compile(slug)
}

// parseString parses a template string.
func (e *Engine) parseString(source string) (*parser.Template, error) {
	l := lexer.New(source)
	p := parser.New(l)
	return p.Parse()
}

// execute executes a parsed template with the given runtime.
func (e *Engine) execute(rt *runtime.Runtime, tmpl *parser.Template) (string, error) {
	err := rt.ExecuteTemplate(tmpl)
	if err != nil {
		return "", err
	}
	return rt.Output(), nil
}

// copyFunctionsToRuntime copies all registered functions to the runtime.
func (e *Engine) copyFunctionsToRuntime(rt *runtime.Runtime) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Copy custom functions (built-ins are already in the runtime)
	// Note: This is a simple implementation. In production, we might want
	// to optimize this by sharing the function registry.
	for name, fn := range e.functions.AllFunctions() {
		rt.RegisterFunction(name, fn)
	}
}

// Config returns a copy of the engine's configuration.
func (e *Engine) Config() Config {
	return e.config
}

// NewWithFS creates an engine from an embedded filesystem.
func NewWithFS(fsys fs.FS, extensions ...string) (*Engine, error) {
	config := DefaultConfig()
	config.TemplateFS = fsys
	if len(extensions) > 0 {
		config.Extensions = extensions
	}
	return New(config)
}

// NewWithDir creates an engine from a directory.
func NewWithDir(dir string, extensions ...string) (*Engine, error) {
	config := DefaultConfig()
	config.TemplateDir = dir
	if len(extensions) > 0 {
		config.Extensions = extensions
	}
	return New(config)
}
