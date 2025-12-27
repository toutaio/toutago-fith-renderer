package loader

import (
	"sync"

	"github.com/toutaio/toutago-fith-renderer/parser"
)

// TemplateCache provides thread-safe caching for parsed templates.
type TemplateCache struct {
	templates map[string]*parser.Template
	mu        sync.RWMutex
}

// NewTemplateCache creates a new template cache.
func NewTemplateCache() *TemplateCache {
	return &TemplateCache{
		templates: make(map[string]*parser.Template),
	}
}

// Get retrieves a template from the cache.
// Returns nil if the template is not cached.
func (c *TemplateCache) Get(slug string) *parser.Template {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.templates[slug]
}

// Set stores a template in the cache.
func (c *TemplateCache) Set(slug string, tmpl *parser.Template) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.templates[slug] = tmpl
}

// Clear removes all templates from the cache.
func (c *TemplateCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.templates = make(map[string]*parser.Template)
}

// Has checks if a template is in the cache.
func (c *TemplateCache) Has(slug string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.templates[slug]
	return exists
}

// Remove removes a template from the cache.
func (c *TemplateCache) Remove(slug string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.templates, slug)
}
