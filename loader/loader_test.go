package loader

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/toutaio/toutago-fith-renderer/parser"
)

//go:embed testdata
var testFS embed.FS

func TestFileSystemLoader_Load(t *testing.T) {
	// Create temp directory with test templates
	tmpDir := t.TempDir()

	// Create test template
	templateContent := "Hello {{.Name}}!"
	templatePath := filepath.Join(tmpDir, "test.html")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Create loader
	loader := NewFileSystemLoader(tmpDir, []string{".html"})

	// Load template
	tmpl, err := loader.Load("test")
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	if tmpl == nil {
		t.Fatal("Expected non-nil template")
	}

	// Verify it was cached
	cached := loader.cache.Get("test")
	if cached == nil {
		t.Error("Expected template to be cached")
	}
}

func TestFileSystemLoader_LoadWithPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested directory structure
	nestedDir := filepath.Join(tmpDir, "layouts")
	err := os.MkdirAll(nestedDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested directory: %v", err)
	}

	// Create template in nested directory
	templateContent := "Layout: {{.Title}}"
	templatePath := filepath.Join(nestedDir, "main.html")
	err = os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	loader := NewFileSystemLoader(tmpDir, []string{".html"})

	// Load using slug
	tmpl, err := loader.Load("layouts/main")
	if err != nil {
		t.Fatalf("Failed to load nested template: %v", err)
	}

	if tmpl == nil {
		t.Fatal("Expected non-nil template")
	}
}

func TestFileSystemLoader_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	loader := NewFileSystemLoader(tmpDir, []string{".html"})

	_, err := loader.Load("nonexistent")
	if err == nil {
		t.Fatal("Expected error for non-existent template")
	}
}

func TestFileSystemLoader_Exists(t *testing.T) {
	tmpDir := t.TempDir()

	templatePath := filepath.Join(tmpDir, "exists.html")
	err := os.WriteFile(templatePath, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	loader := NewFileSystemLoader(tmpDir, []string{".html"})

	if !loader.Exists("exists") {
		t.Error("Expected template to exist")
	}

	if loader.Exists("notexists") {
		t.Error("Expected template to not exist")
	}
}

func TestFileSystemLoader_MultipleExtensions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create templates with different extensions
	err := os.WriteFile(filepath.Join(tmpDir, "test1.html"), []byte("HTML"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "test2.tpl"), []byte("TPL"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	loader := NewFileSystemLoader(tmpDir, []string{".html", ".tpl"})

	// Should find .html first
	tmpl1, err := loader.Load("test1")
	if err != nil {
		t.Fatalf("Failed to load .html template: %v", err)
	}
	if tmpl1 == nil {
		t.Fatal("Expected non-nil template")
	}

	// Should find .tpl
	tmpl2, err := loader.Load("test2")
	if err != nil {
		t.Fatalf("Failed to load .tpl template: %v", err)
	}
	if tmpl2 == nil {
		t.Fatal("Expected non-nil template")
	}
}

func TestFileSystemLoader_ClearCache(t *testing.T) {
	tmpDir := t.TempDir()

	templatePath := filepath.Join(tmpDir, "cached.html")
	err := os.WriteFile(templatePath, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	loader := NewFileSystemLoader(tmpDir, []string{".html"})

	// Load and cache
	_, err = loader.Load("cached")
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	if !loader.cache.Has("cached") {
		t.Fatal("Expected template to be cached")
	}

	// Clear cache
	loader.ClearCache()

	if loader.cache.Has("cached") {
		t.Error("Expected cache to be cleared")
	}
}

func TestEmbedLoader_Load(t *testing.T) {
	loader := NewEmbedLoader(testFS, "testdata", []string{".html"})

	tmpl, err := loader.Load("simple")
	if err != nil {
		t.Fatalf("Failed to load embedded template: %v", err)
	}

	if tmpl == nil {
		t.Fatal("Expected non-nil template")
	}
}

func TestEmbedLoader_Exists(t *testing.T) {
	loader := NewEmbedLoader(testFS, "testdata", []string{".html"})

	if !loader.Exists("simple") {
		t.Error("Expected embedded template to exist")
	}

	if loader.Exists("nonexistent") {
		t.Error("Expected template to not exist")
	}
}

func TestTemplateCache(t *testing.T) {
	cache := NewTemplateCache()

	// Initially empty
	if cache.Has("test") {
		t.Error("Expected cache to be empty")
	}

	if cache.Get("test") != nil {
		t.Error("Expected nil for non-existent template")
	}

	// Set and get
	tmpl := &parser.Template{}
	cache.Set("test", tmpl)

	if !cache.Has("test") {
		t.Error("Expected template to be in cache")
	}

	retrieved := cache.Get("test")
	if retrieved == nil {
		t.Fatal("Expected non-nil template from cache")
	}

	// Remove
	cache.Remove("test")
	if cache.Has("test") {
		t.Error("Expected template to be removed")
	}

	// Clear
	cache.Set("test1", tmpl)
	cache.Set("test2", tmpl)
	cache.Clear()

	if cache.Has("test1") || cache.Has("test2") {
		t.Error("Expected cache to be cleared")
	}
}
