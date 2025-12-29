package compiler

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/toutaio/toutago-fith-renderer/loader"
	"github.com/toutaio/toutago-fith-renderer/parser"
)

// Additional tests for improved coverage

func TestCompiler_ResolveDependenciesWithNested(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create templates
	os.WriteFile(filepath.Join(tmpDir, "main.html"), []byte(`{{include "partial"}}`), 0644)
	os.WriteFile(filepath.Join(tmpDir, "partial.html"), []byte(`{{include "nested"}}`), 0644)
	os.WriteFile(filepath.Join(tmpDir, "nested.html"), []byte("content"), 0644)
	
	ldr := loader.NewFileSystemLoader(tmpDir, []string{".html"})
	c := New(ldr)
	
	tmpl, err := parser.New(`{{include "main"}}`).Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	
	deps, err := c.resolveDependencies(tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(deps) < 1 {
		t.Errorf("expected dependencies, got %d", len(deps))
	}
}

func TestCompiler_ResolveExtendsChain(t *testing.T) {
	tmpDir := t.TempDir()
	
	os.WriteFile(filepath.Join(tmpDir, "base.html"), []byte(`{{block "content"}}base{{end}}`), 0644)
	os.WriteFile(filepath.Join(tmpDir, "middle.html"), []byte(`{{extends "base"}}{{block "content"}}middle{{end}}`), 0644)
	
	ldr := loader.NewFileSystemLoader(tmpDir, []string{".html"})
	c := New(ldr)
	
	tmpl, err := parser.New(`{{extends "middle"}}`).Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	
	deps, err := c.resolveDependencies(tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(deps) < 1 {
		t.Errorf("expected dependencies for extends chain, got %d", len(deps))
	}
}

func TestCompiler_MissingTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	ldr := loader.NewFileSystemLoader(tmpDir, []string{".html"})
	c := New(ldr)
	
	tmpl, err := parser.New(`{{include "missing"}}`).Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	
	_, err = c.resolveDependencies(tmpl)
	if err == nil {
		t.Error("expected error for missing template")
	}
}

func TestCompiler_IfNodeDependencies(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "partial.html"), []byte("content"), 0644)
	
	ldr := loader.NewFileSystemLoader(tmpDir, []string{".html"})
	c := New(ldr)
	
	tmpl, err := parser.New(`{{if .x}}{{include "partial"}}{{end}}`).Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	
	deps, err := c.resolveDependencies(tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(deps) != 1 {
		t.Errorf("expected 1 dependency in if branch, got %d", len(deps))
	}
}

func TestCompiler_RangeNodeDependencies(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "item.html"), []byte("item content"), 0644)
	
	ldr := loader.NewFileSystemLoader(tmpDir, []string{".html"})
	c := New(ldr)
	
	tmpl, err := parser.New(`{{range .items}}{{include "item"}}{{end}}`).Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	
	deps, err := c.resolveDependencies(tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(deps) != 1 {
		t.Errorf("expected 1 dependency in range body, got %d", len(deps))
	}
}

func TestCompiler_BlockNodeDependencies(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "widget.html"), []byte("widget content"), 0644)
	
	ldr := loader.NewFileSystemLoader(tmpDir, []string{".html"})
	c := New(ldr)
	
	tmpl, err := parser.New(`{{block "main"}}{{include "widget"}}{{end}}`).Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	
	deps, err := c.resolveDependencies(tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(deps) != 1 {
		t.Errorf("expected 1 dependency in block, got %d", len(deps))
	}
}

func TestCompiler_DuplicateDependencies(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "common.html"), []byte("common content"), 0644)
	
	ldr := loader.NewFileSystemLoader(tmpDir, []string{".html"})
	c := New(ldr)
	
	// Include same template twice
	tmpl, err := parser.New(`{{include "common"}}{{include "common"}}`).Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	
	deps, err := c.resolveDependencies(tmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Should only list once due to visited check
	if len(deps) != 1 {
		t.Errorf("expected 1 unique dependency, got %d", len(deps))
	}
}

func TestCompiler_CacheKey(t *testing.T) {
	tmpDir := t.TempDir()
	ldr := loader.NewFileSystemLoader(tmpDir, []string{".html"})
	c := New(ldr)
	
	key1 := c.generateCacheKey("test")
	key2 := c.generateCacheKey("test")
	key3 := c.generateCacheKey("other")
	
	if key1 != key2 {
		t.Error("same slug should generate same cache key")
	}
	
	if key1 == key3 {
		t.Error("different slugs should generate different cache keys")
	}
}
