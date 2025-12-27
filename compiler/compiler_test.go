package compiler

import (
	"fmt"
	"testing"

	"github.com/toutaio/toutago-fith-renderer/parser"
)

// mockLoader is a mock loader for testing.
type mockLoader struct {
	templates map[string]*parser.Template
}

func newMockLoader() *mockLoader {
	return &mockLoader{
		templates: make(map[string]*parser.Template),
	}
}

func (m *mockLoader) Load(slug string) (*parser.Template, error) {
	if tmpl, ok := m.templates[slug]; ok {
		return tmpl, nil
	}
	return nil, fmt.Errorf("template %q not found", slug)
}

func (m *mockLoader) Exists(slug string) bool {
	_, ok := m.templates[slug]
	return ok
}

func (m *mockLoader) add(slug string, tmpl *parser.Template) {
	m.templates[slug] = tmpl
}

func TestCompiler_Compile(t *testing.T) {
	loader := newMockLoader()
	loader.add("test", &parser.Template{
		Nodes: []parser.Node{
			&parser.TextNode{Value: "Hello World"},
		},
	})

	compiler := New(loader)
	compiled, err := compiler.Compile("test")
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	if compiled == nil {
		t.Fatal("compiled template is nil")
	}

	if !compiled.IsOptimized {
		t.Error("template should be marked as optimized")
	}

	if len(compiled.AST.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(compiled.AST.Nodes))
	}
}

func TestCompiler_CompileWithCache(t *testing.T) {
	loader := newMockLoader()
	loader.add("cached", &parser.Template{
		Nodes: []parser.Node{
			&parser.TextNode{Value: "Cached template"},
		},
	})

	compiler := New(loader)

	// First compile
	compiled1, err := compiler.Compile("cached")
	if err != nil {
		t.Fatalf("first compile failed: %v", err)
	}

	// Second compile - should hit cache
	compiled2, err := compiler.Compile("cached")
	if err != nil {
		t.Fatalf("second compile failed: %v", err)
	}

	// Should be exact same instance (from cache)
	if compiled1.CacheKey != compiled2.CacheKey {
		t.Error("cache keys should match for same template")
	}
}

func TestCompiler_ResolveDependencies_Include(t *testing.T) {
	loader := newMockLoader()
	loader.add("header", &parser.Template{
		Nodes: []parser.Node{
			&parser.TextNode{Value: "Header"},
		},
	})
	loader.add("main", &parser.Template{
		Nodes: []parser.Node{
			&parser.IncludeNode{Template: "header"},
			&parser.TextNode{Value: "Main content"},
		},
	})

	compiler := New(loader)
	compiled, err := compiler.Compile("main")
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	if len(compiled.Dependencies) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(compiled.Dependencies))
	}

	if len(compiled.Dependencies) > 0 && compiled.Dependencies[0] != "header" {
		t.Errorf("expected dependency 'header', got %q", compiled.Dependencies[0])
	}
}

func TestCompiler_ResolveDependencies_Extends(t *testing.T) {
	loader := newMockLoader()
	loader.add("layout", &parser.Template{
		Nodes: []parser.Node{
			&parser.BlockNode{Name: "content", Body: []parser.Node{
				&parser.TextNode{Value: "Default"},
			}},
		},
	})
	loader.add("page", &parser.Template{
		Nodes: []parser.Node{
			&parser.ExtendsNode{Template: "layout"},
			&parser.BlockNode{Name: "content", Body: []parser.Node{
				&parser.TextNode{Value: "Page content"},
			}},
		},
	})

	compiler := New(loader)
	compiled, err := compiler.Compile("page")
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	if len(compiled.Dependencies) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(compiled.Dependencies))
	}

	if len(compiled.Dependencies) > 0 && compiled.Dependencies[0] != "layout" {
		t.Errorf("expected dependency 'layout', got %q", compiled.Dependencies[0])
	}
}

func TestCompiler_ResolveDependencies_MissingInclude(t *testing.T) {
	loader := newMockLoader()
	loader.add("main", &parser.Template{
		Nodes: []parser.Node{
			&parser.IncludeNode{Template: "missing"},
		},
	})

	compiler := New(loader)
	_, err := compiler.Compile("main")
	if err == nil {
		t.Fatal("expected error for missing include")
	}
}

func TestCompiler_ResolveDependencies_MissingLayout(t *testing.T) {
	loader := newMockLoader()
	loader.add("page", &parser.Template{
		Nodes: []parser.Node{
			&parser.ExtendsNode{Template: "missing-layout"},
		},
	})

	compiler := New(loader)
	_, err := compiler.Compile("page")
	if err == nil {
		t.Fatal("expected error for missing layout")
	}
}

func TestCompiler_ResolveDependencies_Nested(t *testing.T) {
	loader := newMockLoader()
	loader.add("header", &parser.Template{
		Nodes: []parser.Node{
			&parser.TextNode{Value: "Header"},
		},
	})
	loader.add("nav", &parser.Template{
		Nodes: []parser.Node{
			&parser.IncludeNode{Template: "header"},
			&parser.TextNode{Value: "Nav"},
		},
	})
	loader.add("main", &parser.Template{
		Nodes: []parser.Node{
			&parser.IncludeNode{Template: "nav"},
			&parser.TextNode{Value: "Main"},
		},
	})

	compiler := New(loader)
	compiled, err := compiler.Compile("main")
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	// Should have both nav and header as dependencies
	if len(compiled.Dependencies) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(compiled.Dependencies))
	}
}

func TestCompilationCache_GetSet(t *testing.T) {
	cache := NewCompilationCache()

	tmpl := &CompiledTemplate{
		CacheKey:    "test-key",
		IsOptimized: true,
	}

	// Test miss
	if _, ok := cache.Get("test-key"); ok {
		t.Error("should not find template in empty cache")
	}

	// Test set and get
	cache.Set("test-key", tmpl)
	if cached, ok := cache.Get("test-key"); !ok {
		t.Error("should find template after set")
	} else if cached.CacheKey != "test-key" {
		t.Error("cached template has wrong key")
	}
}

func TestCompilationCache_Clear(t *testing.T) {
	cache := NewCompilationCache()

	cache.Set("key1", &CompiledTemplate{CacheKey: "key1"})
	cache.Set("key2", &CompiledTemplate{CacheKey: "key2"})

	cache.Clear()

	if _, ok := cache.Get("key1"); ok {
		t.Error("cache should be empty after clear")
	}
	if _, ok := cache.Get("key2"); ok {
		t.Error("cache should be empty after clear")
	}
}

func TestCompilationCache_Remove(t *testing.T) {
	cache := NewCompilationCache()

	cache.Set("key1", &CompiledTemplate{CacheKey: "key1"})
	cache.Set("key2", &CompiledTemplate{CacheKey: "key2"})

	cache.Remove("key1")

	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 should be removed")
	}
	if _, ok := cache.Get("key2"); !ok {
		t.Error("key2 should still exist")
	}
}

func TestCompiler_ClearCache(t *testing.T) {
	loader := newMockLoader()
	loader.add("test", &parser.Template{
		Nodes: []parser.Node{
			&parser.TextNode{Value: "Test"},
		},
	})

	compiler := New(loader)

	// Compile to populate cache
	_, err := compiler.Compile("test")
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	// Clear cache
	compiler.ClearCache()

	// Cache should be empty (we can't directly test this, but we ensure no panic)
}
