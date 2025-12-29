package compiler

import (
	"testing"

	"github.com/toutaio/toutago-fith-renderer/loader"
	"github.com/toutaio/toutago-fith-renderer/parser"
)

// Additional tests for improved coverage

func TestCompiler_ResolveDependenciesWithNested(t *testing.T) {
	// Create mock loader with nested dependencies
	mockLoader := loader.NewMemoryLoader()
	
	mockLoader.AddTemplate("main", "{{include \"partial\"}}")
	mockLoader.AddTemplate("partial", "{{include \"nested\"}}")
	mockLoader.AddTemplate("nested", "content")
	
	c := New(mockLoader)
	
	tmpl, err := parser.New("{{include \"main\"}}").Parse()
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
	mockLoader := loader.NewMemoryLoader()
	
	mockLoader.AddTemplate("base", "{{block \"content\"}}base{{end}}")
	mockLoader.AddTemplate("middle", "{{extends \"base\"}}{{block \"content\"}}middle{{end}}")
	
	c := New(mockLoader)
	
	tmpl, err := parser.New("{{extends \"middle\"}}").Parse()
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
	mockLoader := loader.NewMemoryLoader()
	c := New(mockLoader)
	
	tmpl, err := parser.New("{{include \"missing\"}}").Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	
	_, err = c.resolveDependencies(tmpl)
	if err == nil {
		t.Error("expected error for missing template")
	}
}

func TestCompiler_IfNodeDependencies(t *testing.T) {
	mockLoader := loader.NewMemoryLoader()
	mockLoader.AddTemplate("partial", "content")
	
	c := New(mockLoader)
	
	tmpl, err := parser.New("{{if .x}}{{include \"partial\"}}{{end}}").Parse()
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
	mockLoader := loader.NewMemoryLoader()
	mockLoader.AddTemplate("item", "item content")
	
	c := New(mockLoader)
	
	tmpl, err := parser.New("{{range .items}}{{include \"item\"}}{{end}}").Parse()
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
	mockLoader := loader.NewMemoryLoader()
	mockLoader.AddTemplate("widget", "widget content")
	
	c := New(mockLoader)
	
	tmpl, err := parser.New("{{block \"main\"}}{{include \"widget\"}}{{end}}").Parse()
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
	mockLoader := loader.NewMemoryLoader()
	mockLoader.AddTemplate("common", "common content")
	
	c := New(mockLoader)
	
	// Include same template twice
	tmpl, err := parser.New("{{include \"common\"}}{{include \"common\"}}").Parse()
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
	c := New(loader.NewMemoryLoader())
	
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
