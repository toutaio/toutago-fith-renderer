package runtime

import (
	"fmt"
	"testing"
	"time"

	"github.com/toutaio/toutago-fith-renderer/lexer"
	"github.com/toutaio/toutago-fith-renderer/parser"
)

// ============================================================================
// String Function Tests
// ============================================================================

func TestFunction_Upper(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "simple upper",
			template: `{{upper "hello"}}`,
			expected: "HELLO",
		},
		{
			name:     "upper with variable",
			template: `{{upper .Name}}`,
			data:     map[string]interface{}{"Name": "alice"},
			expected: "ALICE",
		},
		{
			name:     "upper with pipe",
			template: `{{.Name | upper}}`,
			data:     map[string]interface{}{"Name": "world"},
			expected: "WORLD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeTemplate(tt.template, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if output != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestFunction_Lower(t *testing.T) {
	output, err := executeTemplate(`{{lower "HELLO"}}`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output != "hello" {
		t.Errorf("expected 'hello', got %q", output)
	}
}

func TestFunction_Title(t *testing.T) {
	output, err := executeTemplate(`{{title "hello world"}}`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", output)
	}
}

func TestFunction_Trim(t *testing.T) {
	output, err := executeTemplate(`{{trim "  hello  "}}`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output != "hello" {
		t.Errorf("expected 'hello', got %q", output)
	}
}

func TestFunction_TrimPrefix(t *testing.T) {
	output, err := executeTemplate(`{{trimPrefix "hello world" "hello "}}`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output != "world" {
		t.Errorf("expected 'world', got %q", output)
	}
}

func TestFunction_TrimSuffix(t *testing.T) {
	output, err := executeTemplate(`{{trimSuffix "hello.txt" ".txt"}}`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output != "hello" {
		t.Errorf("expected 'hello', got %q", output)
	}
}

func TestFunction_Truncate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "truncate short string",
			template: `{{truncate "hello" 10}}`,
			expected: "hello",
		},
		{
			name:     "truncate long string",
			template: `{{truncate "hello world this is long" 10}}`,
			expected: "hello worl...",
		},
		{
			name:     "truncate exact length",
			template: `{{truncate "exactly10!" 10}}`,
			expected: "exactly10!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeTemplate(tt.template, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if output != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestFunction_Replace(t *testing.T) {
	output, err := executeTemplate(`{{replace "hello world" "world" "gopher"}}`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output != "hello gopher" {
		t.Errorf("expected 'hello gopher', got %q", output)
	}
}

// ============================================================================
// Array Function Tests
// ============================================================================

func TestFunction_Join(t *testing.T) {
	data := map[string]interface{}{
		"Items": []string{"apple", "banana", "cherry"},
	}

	output, err := executeTemplate(`{{join .Items ", "}}`, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "apple, banana, cherry"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestFunction_Len(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     interface{}
		expected string
	}{
		{
			name:     "len of array",
			template: `{{len .Items}}`,
			data:     map[string]interface{}{"Items": []int{1, 2, 3, 4}},
			expected: "4",
		},
		{
			name:     "len of string",
			template: `{{len "hello"}}`,
			expected: "5",
		},
		{
			name:     "len of map",
			template: `{{len .Config}}`,
			data:     map[string]interface{}{"Config": map[string]int{"a": 1, "b": 2}},
			expected: "2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeTemplate(tt.template, tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if output != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestFunction_First(t *testing.T) {
	data := map[string]interface{}{
		"Items": []string{"first", "second", "third"},
	}

	output, err := executeTemplate(`{{first .Items}}`, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "first" {
		t.Errorf("expected 'first', got %q", output)
	}
}

func TestFunction_Last(t *testing.T) {
	data := map[string]interface{}{
		"Items": []string{"first", "second", "third"},
	}

	output, err := executeTemplate(`{{last .Items}}`, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "third" {
		t.Errorf("expected 'third', got %q", output)
	}
}

func TestFunction_FirstEmpty(t *testing.T) {
	data := map[string]interface{}{
		"Items": []string{},
	}

	_, err := executeTemplate(`{{first .Items}}`, data)
	if err == nil {
		t.Error("expected error for empty array")
	}
}

// ============================================================================
// Logic Function Tests
// ============================================================================

func TestFunction_Default(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     interface{}
		expected string
	}{
		{
			name:     "default with empty string",
			template: `{{default .Name "Anonymous"}}`,
			data:     map[string]interface{}{"Name": ""},
			expected: "Anonymous",
		},
		{
			name:     "default with non-empty string",
			template: `{{default .Name "Anonymous"}}`,
			data:     map[string]interface{}{"Name": "Alice"},
			expected: "Alice",
		},
		{
			name:     "default with zero",
			template: `{{default .Count 10}}`,
			data:     map[string]interface{}{"Count": 0},
			expected: "10",
		},
		{
			name:     "default with non-zero",
			template: `{{default .Count 10}}`,
			data:     map[string]interface{}{"Count": 5},
			expected: "5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeTemplate(tt.template, tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if output != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, output)
			}
		})
	}
}

// ============================================================================
// Encoding Function Tests
// ============================================================================

func TestFunction_URLEncode(t *testing.T) {
	output, err := executeTemplate(`{{urlEncode "hello world"}}`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "hello+world"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestFunction_HTMLEscape(t *testing.T) {
	output, err := executeTemplate(`{{htmlEscape "<script>alert('xss')</script>"}}`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

// ============================================================================
// Date Function Tests
// ============================================================================

func TestFunction_Date(t *testing.T) {
	// Create a specific time
	testTime := time.Date(2024, 12, 25, 15, 30, 45, 0, time.UTC)

	data := map[string]interface{}{
		"Time": testTime,
	}

	output, err := executeTemplate(`{{date "YYYY-MM-DD" .Time}}`, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "2024-12-25"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestFunction_DateWithString(t *testing.T) {
	data := map[string]interface{}{
		"DateStr": "2024-12-25",
	}

	output, err := executeTemplate(`{{date "YYYY-MM-DD" .DateStr}}`, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "2024-12-25"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

// ============================================================================
// Pipe Chain Tests
// ============================================================================

func TestFunction_PipeChain(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     interface{}
		expected string
	}{
		{
			name:     "double pipe",
			template: `{{.Name | lower | trim}}`,
			data:     map[string]interface{}{"Name": "  HELLO  "},
			expected: "hello",
		},
		{
			name:     "triple pipe",
			template: `{{.Text | trim | upper | htmlEscape}}`,
			data:     map[string]interface{}{"Text": "  <hello>  "},
			expected: "&lt;HELLO&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeTemplate(tt.template, tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if output != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, output)
			}
		})
	}
}

// ============================================================================
// Complex Integration Tests
// ============================================================================

func TestFunction_ComplexTemplate(t *testing.T) {
	data := map[string]interface{}{
		"Title": "my article",
		"Tags":  []string{"go", "templates", "web"},
		"Count": 0,
	}

	template := `{{title .Title | upper}}
Tags: {{join .Tags ", " | upper}}
Items: {{default .Count 5}}`

	output, err := executeTemplate(template, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `MY ARTICLE
Tags: GO, TEMPLATES, WEB
Items: 5`

	if output != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, output)
	}
}

// ============================================================================
// Custom Function Tests
// ============================================================================

func TestRuntime_CustomFunction(t *testing.T) {
	ctx := NewContext(map[string]interface{}{
		"Value": 5,
	})

	rt := NewRuntime(ctx)

	// Register custom function
	rt.RegisterFunction("double", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("double: expected 1 argument")
		}
		n, ok := args[0].(int)
		if !ok {
			return nil, fmt.Errorf("double: argument must be an integer")
		}
		return n * 2, nil
	})

	// Parse template
	l := lexer.New("{{double .Value}}")
	p := parser.New(l)
	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// Execute with the runtime that has the custom function
	for _, node := range ast.Nodes {
		if err := rt.executeNode(node); err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	output := rt.output.String()
	if output != "10" {
		t.Errorf("expected '10', got %q", output)
	}
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestFunction_ErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     interface{}
	}{
		{
			name:     "unknown function",
			template: `{{unknownFunc "test"}}`,
		},
		{
			name:     "wrong arg count",
			template: `{{upper "a" "b"}}`,
		},
		{
			name:     "wrong arg type",
			template: `{{upper 123}}`,
		},
		{
			name:     "first on empty array",
			template: `{{first .Items}}`,
			data:     map[string]interface{}{"Items": []string{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executeTemplate(tt.template, tt.data)
			if err == nil {
				t.Error("expected error but got none")
			}
		})
	}
}
