// Package fith provides configuration for the Fíth template engine.
package fith

import (
	"io/fs"
)

// Default configuration constants
const (
	defaultTemplateDir    = "templates"
	defaultLeftDelimiter  = "{{"
	defaultRightDelimiter = "}}"
)

// Config configures the Fíth template engine.
type Config struct {
	// TemplateDir is the base directory for template files.
	// Used when loading templates from the filesystem.
	TemplateDir string

	// TemplateFS is an embedded filesystem for templates.
	// If set, takes precedence over TemplateDir.
	TemplateFS fs.FS

	// Extensions are file extensions to try when resolving template slugs.
	// Default: [".html", ".tpl", ".txt"]
	Extensions []string

	// LeftDelimiter is the opening delimiter for template expressions.
	// Default: "{{"
	LeftDelimiter string

	// RightDelimiter is the closing delimiter for template expressions.
	// Default: "}}"
	RightDelimiter string

	// CacheEnabled enables template compilation caching.
	// Default: true
	CacheEnabled bool

	// AutoEscape enables automatic HTML escaping of variables.
	// Default: false (manual escaping via htmlEscape function)
	AutoEscape bool

	// StrictMode causes the engine to fail on undefined variables.
	// Default: false (undefined variables render as empty string)
	StrictMode bool

	// MaxIncludeDepth limits the depth of template includes to prevent infinite recursion.
	// Default: 100
	MaxIncludeDepth int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		TemplateDir:     "templates",
		Extensions:      []string{".html", ".tpl", ".txt"},
		LeftDelimiter:   "{{",
		RightDelimiter:  "}}",
		CacheEnabled:    true,
		AutoEscape:      false,
		StrictMode:      false,
		MaxIncludeDepth: 100,
	}
}

// Validate checks the configuration for errors.
func (c *Config) Validate() error {
	if c.TemplateDir == "" && c.TemplateFS == nil {
		return NewError(ErrorTypeLoader, "either TemplateDir or TemplateFS must be set")
	}

	if c.LeftDelimiter == "" || c.RightDelimiter == "" {
		return NewError(ErrorTypeTemplate, "delimiters cannot be empty")
	}

	if c.LeftDelimiter == c.RightDelimiter {
		return NewError(ErrorTypeTemplate, "left and right delimiters must be different")
	}

	if c.MaxIncludeDepth < 1 {
		return NewError(ErrorTypeTemplate, "MaxIncludeDepth must be at least 1")
	}

	if len(c.Extensions) == 0 {
		c.Extensions = []string{".html", ".tpl", ".txt"}
	}

	return nil
}

// applyDefaults applies default values to missing configuration options.
func (c *Config) applyDefaults() {
	if c.TemplateDir == "" && c.TemplateFS == nil {
		c.TemplateDir = defaultTemplateDir
	}

	if len(c.Extensions) == 0 {
		c.Extensions = []string{".html", ".tpl", ".txt"}
	}

	if c.LeftDelimiter == "" {
		c.LeftDelimiter = defaultLeftDelimiter
	}

	if c.RightDelimiter == "" {
		c.RightDelimiter = defaultRightDelimiter
	}

	if c.MaxIncludeDepth == 0 {
		c.MaxIncludeDepth = 100
	}
}
