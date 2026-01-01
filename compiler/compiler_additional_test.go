package compiler

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/toutaio/toutago-fith-renderer/lexer"
	"github.com/toutaio/toutago-fith-renderer/loader"
	"github.com/toutaio/toutago-fith-renderer/parser"
)

func TestNewCompiler(t *testing.T) {
	tmpDir := t.TempDir()
	ldr := loader.NewFileSystemLoader(tmpDir, []string{".html"})
	c := NewCompiler(ldr)
	if c == nil {
		t.Fatal("NewCompiler returned nil")
	}
}

func TestCompileWithoutCache(t *testing.T) {
	tmpDir := t.TempDir()
	content := "Hello {{.name}}"
	_ = os.WriteFile(filepath.Join(tmpDir, "test.html"), []byte(content), 0o600)

	ldr := loader.NewFileSystemLoader(tmpDir, []string{".html"})
	c := New(ldr)

	// Parse template
	l := lexer.New(content)
	p := parser.New(l)
	parsedTmpl, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Test CompileWithoutCache
	compiled, err := c.CompileWithoutCache(parsedTmpl)
	if err != nil {
		t.Fatalf("CompileWithoutCache failed: %v", err)
	}
	if compiled == nil {
		t.Fatal("CompileWithoutCache returned nil")
	}
}

func TestResolveIfNodeDeps(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a template with if node containing includes
	_ = os.WriteFile(filepath.Join(tmpDir, "main.html"), []byte(`{{if .show}}{{include "part"}}{{end}}`), 0o600)
	_ = os.WriteFile(filepath.Join(tmpDir, "part.html"), []byte("Partial"), 0o600)

	ldr := loader.NewFileSystemLoader(tmpDir, []string{".html"})
	c := New(ldr)

	_, err := c.Compile("main")
	if err != nil {
		t.Fatalf("Compile with if node deps failed: %v", err)
	}
}
