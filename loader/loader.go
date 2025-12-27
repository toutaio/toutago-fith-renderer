// Package loader provides template loading capabilities for the Fíth template engine.
// It supports loading templates from various sources including filesystems and embedded files.
package loader

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/toutaio/toutago-fith-renderer/lexer"
	"github.com/toutaio/toutago-fith-renderer/parser"
)

// Loader is the interface for loading templates from various sources.
// Implementations can load from filesystem, embedded FS, memory, etc.
type Loader interface {
	// Load loads a template by its slug (e.g., "layouts/main", "partials/header").
	// Returns the parsed template or an error if not found or parsing fails.
	Load(slug string) (*parser.Template, error)

	// Exists checks if a template exists without loading it.
	Exists(slug string) bool
}

// FileSystemLoader loads templates from a directory on the filesystem.
type FileSystemLoader struct {
	baseDir    string
	extensions []string
	cache      *TemplateCache
	mu         sync.RWMutex
}

// NewFileSystemLoader creates a new filesystem-based template loader.
// baseDir is the root directory for templates.
// extensions are file extensions to try (e.g., ".html", ".tpl").
func NewFileSystemLoader(baseDir string, extensions []string) *FileSystemLoader {
	if len(extensions) == 0 {
		extensions = []string{".html", ".tpl", ".txt"}
	}
	return &FileSystemLoader{
		baseDir:    baseDir,
		extensions: extensions,
		cache:      NewTemplateCache(),
	}
}

// Load loads and parses a template by slug.
// Slug format: "layouts/main" resolves to "baseDir/layouts/main.{ext}"
func (l *FileSystemLoader) Load(slug string) (*parser.Template, error) {
	// Check cache first
	if tmpl := l.cache.Get(slug); tmpl != nil {
		return tmpl, nil
	}

	// Find the template file
	path, err := l.resolvePath(slug)
	if err != nil {
		return nil, err
	}

	// Read and parse
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %q: %w", slug, err)
	}

	tmpl, err := l.parse(string(content), slug)
	if err != nil {
		return nil, err
	}

	// Cache the parsed template
	l.cache.Set(slug, tmpl)
	return tmpl, nil
}

// Exists checks if a template file exists.
func (l *FileSystemLoader) Exists(slug string) bool {
	path, err := l.resolvePath(slug)
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

// resolvePath resolves a slug to a filesystem path.
// Tries each extension until a file is found.
func (l *FileSystemLoader) resolvePath(slug string) (string, error) {
	// Convert slug separators to OS path separators
	slug = filepath.FromSlash(slug)

	for _, ext := range l.extensions {
		// Try with extension
		path := filepath.Join(l.baseDir, slug+ext)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

		// Try without adding extension if slug already has it
		if strings.HasSuffix(slug, ext) {
			path := filepath.Join(l.baseDir, slug)
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}
	}

	return "", fmt.Errorf("template %q not found in %q", slug, l.baseDir)
}

// parse parses template content into an AST.
func (l *FileSystemLoader) parse(content, slug string) (*parser.Template, error) {
	lex := lexer.New(content)
	p := parser.New(lex)
	tmpl, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %q: %w", slug, err)
	}
	return tmpl, nil
}

// ClearCache clears the template cache.
func (l *FileSystemLoader) ClearCache() {
	l.cache.Clear()
}

// EmbedLoader loads templates from an embedded filesystem (embed.FS).
type EmbedLoader struct {
	fs         fs.FS
	baseDir    string
	extensions []string
	cache      *TemplateCache
}

// NewEmbedLoader creates a new embed.FS-based template loader.
func NewEmbedLoader(fsys fs.FS, baseDir string, extensions []string) *EmbedLoader {
	if len(extensions) == 0 {
		extensions = []string{".html", ".tpl", ".txt"}
	}
	return &EmbedLoader{
		fs:         fsys,
		baseDir:    baseDir,
		extensions: extensions,
		cache:      NewTemplateCache(),
	}
}

// Load loads and parses a template by slug from embedded filesystem.
func (l *EmbedLoader) Load(slug string) (*parser.Template, error) {
	// Check cache first
	if tmpl := l.cache.Get(slug); tmpl != nil {
		return tmpl, nil
	}

	// Find the template file
	path, err := l.resolvePath(slug)
	if err != nil {
		return nil, err
	}

	// Read and parse
	content, err := fs.ReadFile(l.fs, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %q: %w", slug, err)
	}

	tmpl, err := l.parse(string(content), slug)
	if err != nil {
		return nil, err
	}

	// Cache the parsed template
	l.cache.Set(slug, tmpl)
	return tmpl, nil
}

// Exists checks if a template file exists in embedded filesystem.
func (l *EmbedLoader) Exists(slug string) bool {
	path, err := l.resolvePath(slug)
	if err != nil {
		return false
	}
	_, err = fs.Stat(l.fs, path)
	return err == nil
}

// resolvePath resolves a slug to an embedded filesystem path.
func (l *EmbedLoader) resolvePath(slug string) (string, error) {
	// Use forward slashes for embed.FS
	slug = filepath.ToSlash(slug)

	for _, ext := range l.extensions {
		// Try with extension
		path := filepath.Join(l.baseDir, slug+ext)
		path = filepath.ToSlash(path) // Ensure forward slashes
		if _, err := fs.Stat(l.fs, path); err == nil {
			return path, nil
		}

		// Try without adding extension if slug already has it
		if strings.HasSuffix(slug, ext) {
			path := filepath.Join(l.baseDir, slug)
			path = filepath.ToSlash(path)
			if _, err := fs.Stat(l.fs, path); err == nil {
				return path, nil
			}
		}
	}

	return "", fmt.Errorf("template %q not found in embedded filesystem", slug)
}

// parse parses template content into an AST.
func (l *EmbedLoader) parse(content, slug string) (*parser.Template, error) {
	lex := lexer.New(content)
	p := parser.New(lex)
	tmpl, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %q: %w", slug, err)
	}
	return tmpl, nil
}

// ClearCache clears the template cache.
func (l *EmbedLoader) ClearCache() {
	l.cache.Clear()
}
