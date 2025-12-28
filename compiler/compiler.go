package compiler

import (
	"fmt"
	"hash/fnv"
	"sync"

	"github.com/toutaio/toutago-fith-renderer/parser"
)

// Compiler compiles and optimizes templates.
type Compiler struct {
	cache     *CompilationCache
	loader    TemplateLoader
	optimizer *Optimizer
}

// TemplateLoader defines the interface for loading template source.
type TemplateLoader interface {
	Load(slug string) (*parser.Template, error)
	Exists(slug string) bool
}

// CompiledTemplate represents an optimized, executable template.
type CompiledTemplate struct {
	AST          *parser.Template
	Dependencies []string
	CacheKey     string
	IsOptimized  bool
}

// CompilationCache provides thread-safe caching of compiled templates.
type CompilationCache struct {
	mu        sync.RWMutex
	templates map[string]*CompiledTemplate
}

// NewCompilationCache creates a new compilation cache.
func NewCompilationCache() *CompilationCache {
	return &CompilationCache{
		templates: make(map[string]*CompiledTemplate),
	}
}

// Get retrieves a compiled template from cache.
func (c *CompilationCache) Get(key string) (*CompiledTemplate, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	tmpl, ok := c.templates[key]
	return tmpl, ok
}

// Set stores a compiled template in cache.
func (c *CompilationCache) Set(key string, tmpl *CompiledTemplate) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.templates[key] = tmpl
}

// Clear removes all cached templates.
func (c *CompilationCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.templates = make(map[string]*CompiledTemplate)
}

// Remove removes a specific template from cache.
func (c *CompilationCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.templates, key)
}

// New creates a new compiler with the given loader.
func New(loader TemplateLoader) *Compiler {
	return &Compiler{
		cache:     NewCompilationCache(),
		loader:    loader,
		optimizer: NewOptimizer(),
	}
}

// NewCompiler is an alias for New.
func NewCompiler(loader TemplateLoader) *Compiler {
	return New(loader)
}

// Compile compiles a template by slug, with caching and optimization.
func (c *Compiler) Compile(slug string) (*CompiledTemplate, error) {
	// Generate cache key
	cacheKey := c.generateCacheKey(slug)

	// Check cache
	if cached, ok := c.cache.Get(cacheKey); ok {
		return cached, nil
	}

	// Load template
	tmpl, err := c.loader.Load(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to load template %q: %w", slug, err)
	}

	// Resolve dependencies
	deps, err := c.resolveDependencies(tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies for %q: %w", slug, err)
	}

	// Optimize AST
	optimized := c.optimizer.Optimize(tmpl)

	// Create compiled template
	compiled := &CompiledTemplate{
		AST:          optimized,
		Dependencies: deps,
		CacheKey:     cacheKey,
		IsOptimized:  true,
	}

	// Cache it
	c.cache.Set(cacheKey, compiled)

	return compiled, nil
}

// CompileWithoutCache compiles a template without using the cache.
func (c *Compiler) CompileWithoutCache(tmpl *parser.Template) (*CompiledTemplate, error) {
	// Resolve dependencies
	deps, err := c.resolveDependencies(tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	// Optimize AST
	optimized := c.optimizer.Optimize(tmpl)

	// Create compiled template
	compiled := &CompiledTemplate{
		AST:          optimized,
		Dependencies: deps,
		CacheKey:     "",
		IsOptimized:  true,
	}

	return compiled, nil
}

// ClearCache clears the compilation cache.
func (c *Compiler) ClearCache() {
	c.cache.Clear()
}

// generateCacheKey generates a unique cache key for a template slug.
func (c *Compiler) generateCacheKey(slug string) string {
	h := fnv.New64a()
	h.Write([]byte(slug))
	return fmt.Sprintf("%s-%x", slug, h.Sum64())
}

// resolveDependencies finds all template dependencies (includes, extends).
func (c *Compiler) resolveDependencies(tmpl *parser.Template) ([]string, error) {
	deps := make([]string, 0)
	visited := make(map[string]bool)

	var resolve func(*parser.Template) error
	resolve = func(t *parser.Template) error {
		for _, node := range t.Nodes {
			switch n := node.(type) {
			case *parser.IncludeNode:
				if !visited[n.Template] {
					visited[n.Template] = true
					deps = append(deps, n.Template)

					// Verify include exists
					if !c.loader.Exists(n.Template) {
						return fmt.Errorf("included template %q not found", n.Template)
					}

					// Load and recursively resolve
					included, err := c.loader.Load(n.Template)
					if err != nil {
						return fmt.Errorf("failed to load included template %q: %w", n.Template, err)
					}
					if err := resolve(included); err != nil {
						return err
					}
				}

			case *parser.ExtendsNode:
				if !visited[n.Template] {
					visited[n.Template] = true
					deps = append(deps, n.Template)

					// Verify layout exists
					if !c.loader.Exists(n.Template) {
						return fmt.Errorf("layout template %q not found", n.Template)
					}

					// Load and recursively resolve
					layout, err := c.loader.Load(n.Template)
					if err != nil {
						return fmt.Errorf("failed to load layout %q: %w", n.Template, err)
					}
					if err := resolve(layout); err != nil {
						return err
					}
				}

			case *parser.IfNode:
				// Check then and else branches
				for _, child := range n.Then {
					if err := resolve(&parser.Template{Nodes: []parser.Node{child}}); err != nil {
						return err
					}
				}
				if n.Else != nil {
					for _, child := range n.Else {
						if err := resolve(&parser.Template{Nodes: []parser.Node{child}}); err != nil {
							return err
						}
					}
				}

			case *parser.RangeNode:
				// Check loop body
				for _, child := range n.Body {
					if err := resolve(&parser.Template{Nodes: []parser.Node{child}}); err != nil {
						return err
					}
				}

			case *parser.BlockNode:
				// Check block body
				for _, child := range n.Body {
					if err := resolve(&parser.Template{Nodes: []parser.Node{child}}); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}

	if err := resolve(tmpl); err != nil {
		return nil, err
	}

	return deps, nil
}
