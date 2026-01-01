package runtime

import (
	"strings"
	"testing"

	"github.com/toutaio/toutago-fith-renderer/lexer"
	"github.com/toutaio/toutago-fith-renderer/parser"
)

func TestAllFunctions(t *testing.T) {
	reg := NewFunctionRegistry()
	funcs := reg.AllFunctions()

	if len(funcs) == 0 {
		t.Error("AllFunctions returned empty map")
	}

	// Check for some known functions
	if _, ok := funcs["upper"]; !ok {
		t.Error("AllFunctions missing 'upper' function")
	}
	if _, ok := funcs["lower"]; !ok {
		t.Error("AllFunctions missing 'lower' function")
	}
}

func TestRuntimeGetContext(t *testing.T) {
	ctx := NewContext(map[string]interface{}{"key": "value"})
	rt := NewRuntime(ctx)

	gotCtx := rt.GetContext()
	if gotCtx != ctx {
		t.Error("GetContext returned different context")
	}
}

func TestRuntimeExecuteTemplateAndOutput(t *testing.T) {
	l := lexer.New("Hello {{.name}}")
	p := parser.New(l)
	tmpl, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	ctx := NewContext(map[string]interface{}{"name": "World"})
	rt := NewRuntime(ctx)

	err = rt.ExecuteTemplate(tmpl)
	if err != nil {
		t.Fatalf("ExecuteTemplate failed: %v", err)
	}

	output := rt.Output()
	if !strings.Contains(output, "Hello World") {
		t.Errorf("Expected 'Hello World', got: %s", output)
	}
}
