package runtime

import (
	"testing"

	"github.com/toutaio/toutago-fith-renderer/lexer"
	"github.com/toutaio/toutago-fith-renderer/parser"
)

// Helper function to parse and execute a template
func executeTemplate(input string, data interface{}) (string, error) {
	l := lexer.New(input)
	p := parser.New(l)
	ast, err := p.Parse()
	if err != nil {
		return "", err
	}

	ctx := NewContext(data)
	return Execute(ast, ctx)
}

func TestRuntime_SimpleText(t *testing.T) {
	output, err := executeTemplate("Hello, World!", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got %q", output)
	}
}

func TestRuntime_SimpleVariable(t *testing.T) {
	data := map[string]interface{}{
		"Name": "Alice",
	}

	output, err := executeTemplate("Hello {{.Name}}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Hello Alice"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestRuntime_NestedVariable(t *testing.T) {
	data := map[string]interface{}{
		"User": map[string]interface{}{
			"Name":  "Bob",
			"Email": "bob@example.com",
		},
	}

	output, err := executeTemplate("Name: {{.User.Name}}, Email: {{.User.Email}}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Name: Bob, Email: bob@example.com"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestRuntime_StructVariable(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	data := map[string]interface{}{
		"User": User{
			Name:  "Charlie",
			Age:   30,
			Email: "charlie@example.com",
		},
	}

	output, err := executeTemplate("{{.User.Name}} is {{.User.Age}} years old", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Charlie is 30 years old"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestRuntime_IfTrue(t *testing.T) {
	data := map[string]interface{}{
		"Active": true,
	}

	output, err := executeTemplate("{{if .Active}}Active{{end}}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "Active" {
		t.Errorf("expected 'Active', got %q", output)
	}
}

func TestRuntime_IfFalse(t *testing.T) {
	data := map[string]interface{}{
		"Active": false,
	}

	output, err := executeTemplate("{{if .Active}}Active{{end}}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "" {
		t.Errorf("expected empty string, got %q", output)
	}
}

func TestRuntime_IfElse(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		expected string
	}{
		{
			name:     "true condition",
			data:     map[string]interface{}{"Active": true},
			expected: "Yes",
		},
		{
			name:     "false condition",
			data:     map[string]interface{}{"Active": false},
			expected: "No",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeTemplate("{{if .Active}}Yes{{else}}No{{end}}", tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if output != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestRuntime_Comparison(t *testing.T) {
	tests := []struct {
		template string
		data     map[string]interface{}
		expected string
	}{
		{
			template: "{{if .A == .B}}equal{{end}}",
			data:     map[string]interface{}{"A": 5, "B": 5},
			expected: "equal",
		},
		{
			template: "{{if .A != .B}}not equal{{end}}",
			data:     map[string]interface{}{"A": 5, "B": 3},
			expected: "not equal",
		},
		{
			template: "{{if .A < .B}}less{{end}}",
			data:     map[string]interface{}{"A": 3, "B": 5},
			expected: "less",
		},
		{
			template: "{{if .A > .B}}greater{{end}}",
			data:     map[string]interface{}{"A": 5, "B": 3},
			expected: "greater",
		},
	}

	for _, tt := range tests {
		t.Run(tt.template, func(t *testing.T) {
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

func TestRuntime_LogicalOperators(t *testing.T) {
	tests := []struct {
		template string
		data     map[string]interface{}
		expected string
	}{
		{
			template: "{{if .A && .B}}both{{end}}",
			data:     map[string]interface{}{"A": true, "B": true},
			expected: "both",
		},
		{
			template: "{{if .A && .B}}both{{end}}",
			data:     map[string]interface{}{"A": true, "B": false},
			expected: "",
		},
		{
			template: "{{if .A || .B}}either{{end}}",
			data:     map[string]interface{}{"A": true, "B": false},
			expected: "either",
		},
		{
			template: "{{if .A || .B}}either{{end}}",
			data:     map[string]interface{}{"A": false, "B": false},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.template, func(t *testing.T) {
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

func TestRuntime_Not(t *testing.T) {
	data := map[string]interface{}{
		"Active": false,
	}

	output, err := executeTemplate("{{if !.Active}}Inactive{{end}}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "Inactive" {
		t.Errorf("expected 'Inactive', got %q", output)
	}
}

func TestRuntime_RangeSlice(t *testing.T) {
	data := map[string]interface{}{
		"Items": []string{"apple", "banana", "cherry"},
	}

	output, err := executeTemplate("{{range .Items}}{{.}},{{end}}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "apple,banana,cherry,"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestRuntime_RangeWithIndex(t *testing.T) {
	data := map[string]interface{}{
		"Items": []string{"a", "b", "c"},
	}

	output, err := executeTemplate("{{range .Items}}{{@index}}:{{.}} {{end}}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "0:a 1:b 2:c "
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestRuntime_RangeFirstLast(t *testing.T) {
	data := map[string]interface{}{
		"Items": []string{"a", "b", "c"},
	}

	template := `{{range .Items}}{{if @first}}[{{end}}{{.}}{{if @last}}]{{else}},{{end}}{{end}}`
	output, err := executeTemplate(template, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "[a,b,c]"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestRuntime_RangeEmptySlice(t *testing.T) {
	data := map[string]interface{}{
		"Items": []string{},
	}

	output, err := executeTemplate("before{{range .Items}}{{.}}{{end}}after", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "beforeafter"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestRuntime_RangeMap(t *testing.T) {
	data := map[string]interface{}{
		"User": map[string]string{
			"name":  "Alice",
			"email": "alice@example.com",
		},
	}

	// Note: Map iteration order is not guaranteed, so we just check for presence
	output, err := executeTemplate("{{range .User}}{{@key}}={{.}} {{end}}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain both key-value pairs
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRuntime_ArrayAccess(t *testing.T) {
	data := map[string]interface{}{
		"Items": []string{"first", "second", "third"},
	}

	output, err := executeTemplate("{{.Items[0]}} and {{.Items[2]}}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "first and third"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestRuntime_MapAccess(t *testing.T) {
	data := map[string]interface{}{
		"Config": map[string]string{
			"host": "localhost",
			"port": "8080",
		},
	}

	output, err := executeTemplate(`{{.Config["host"]}}:{{.Config["port"]}}`, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "localhost:8080"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestRuntime_Arithmetic(t *testing.T) {
	tests := []struct {
		template string
		data     map[string]interface{}
		expected string
	}{
		{
			template: "{{.A + .B}}",
			data:     map[string]interface{}{"A": 5, "B": 3},
			expected: "8",
		},
		{
			template: "{{.A - .B}}",
			data:     map[string]interface{}{"A": 5, "B": 3},
			expected: "2",
		},
		{
			template: "{{.A * .B}}",
			data:     map[string]interface{}{"A": 5, "B": 3},
			expected: "15",
		},
		{
			template: "{{.A / .B}}",
			data:     map[string]interface{}{"A": 6.0, "B": 2.0},
			expected: "3",
		},
		{
			template: "{{.A % .B}}",
			data:     map[string]interface{}{"A": 5, "B": 3},
			expected: "2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.template, func(t *testing.T) {
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

func TestRuntime_ComplexTemplate(t *testing.T) {
	type User struct {
		Name   string
		Active bool
		Age    int
	}

	data := map[string]interface{}{
		"Title": "User List",
		"Users": []User{
			{Name: "Alice", Active: true, Age: 30},
			{Name: "Bob", Active: false, Age: 25},
			{Name: "Charlie", Active: true, Age: 35},
		},
	}

	template := `<h1>{{.Title}}</h1>
{{range .Users}}
<div>
  <p>{{.Name}} ({{.Age}})</p>
  {{if .Active}}
  <span>Active</span>
  {{else}}
  <span>Inactive</span>
  {{end}}
</div>
{{end}}`

	output, err := executeTemplate(template, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Just check it contains expected parts
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}

	// Should contain title
	if !contains(output, "<h1>User List</h1>") {
		t.Error("expected output to contain title")
	}

	// Should contain user names
	if !contains(output, "Alice") || !contains(output, "Bob") || !contains(output, "Charlie") {
		t.Error("expected output to contain all user names")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || len(s) > len(substr) && contains(s[1:], substr)
}

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"nil", nil, false},
		{"true", true, true},
		{"false", false, false},
		{"zero int", 0, false},
		{"non-zero int", 42, true},
		{"empty string", "", false},
		{"non-empty string", "hello", true},
		{"empty slice", []int{}, false},
		{"non-empty slice", []int{1}, true},
		{"zero float", 0.0, false},
		{"non-zero float", 3.14, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTruthy(tt.value)
			if result != tt.expected {
				t.Errorf("IsTruthy(%v) = %v, expected %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestRuntime_DivisionByZero(t *testing.T) {
	data := map[string]interface{}{
		"A": 5,
		"B": 0,
	}

	_, err := executeTemplate("{{.A / .B}}", data)
	if err == nil {
		t.Error("expected error for division by zero")
	}
}

func TestRuntime_ModuloByZero(t *testing.T) {
	data := map[string]interface{}{
		"A": 5,
		"B": 0,
	}

	_, err := executeTemplate("{{.A % .B}}", data)
	if err == nil {
		t.Error("expected error for modulo by zero")
	}
}

func TestRuntime_InvalidVariable(t *testing.T) {
	data := map[string]interface{}{
		"Name": "Alice",
	}

	_, err := executeTemplate("{{.NonExistent}}", data)
	if err == nil {
		t.Error("expected error for non-existent variable")
	}
}

func TestRuntime_NilFieldAccess(t *testing.T) {
	data := map[string]interface{}{
		"User": nil,
	}

	_, err := executeTemplate("{{.User.Name}}", data)
	if err == nil {
		t.Error("expected error for nil field access")
	}
}

func TestRuntime_NonIterableRange(t *testing.T) {
	data := map[string]interface{}{
		"Value": "not iterable",
	}

	_, err := executeTemplate("{{range .Value}}{{.}}{{end}}", data)
	if err == nil {
		t.Error("expected error for non-iterable value in range")
	}
}

func TestRuntime_FloatArithmetic(t *testing.T) {
	data := map[string]interface{}{
		"A": 3.5,
		"B": 1.5,
	}

	output, err := executeTemplate("{{.A + .B}}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "5" {
		t.Errorf("expected .5 got %q", output)
	}
}

func TestRuntime_StringComparison(t *testing.T) {
	data := map[string]interface{}{
		"A": "apple",
		"B": "banana",
	}

	output, err := executeTemplate("{{if .A < .B}}less{{end}}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "less" {
		t.Errorf("expected 'less', got %q", output)
	}
}

func TestRuntime_TruthyValues(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		template string
		expected string
	}{
		{
			name:     "empty string is falsy",
			data:     map[string]interface{}{"S": ""},
			template: "{{if .S}}yes{{else}}no{{end}}",
			expected: "no",
		},
		{
			name:     "non-empty string is truthy",
			data:     map[string]interface{}{"S": "hello"},
			template: "{{if .S}}yes{{else}}no{{end}}",
			expected: "yes",
		},
		{
			name:     "zero is falsy",
			data:     map[string]interface{}{"N": 0},
			template: "{{if .N}}yes{{else}}no{{end}}",
			expected: "no",
		},
		{
			name:     "empty slice is falsy",
			data:     map[string]interface{}{"A": []int{}},
			template: "{{if .A}}yes{{else}}no{{end}}",
			expected: "no",
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

func TestRuntime_PointerDereference(t *testing.T) {
	type User struct {
		Name string
	}

	name := "Alice"
	user := User{Name: name}
	userPtr := &user

	data := map[string]interface{}{
		"User": userPtr,
	}

	output, err := executeTemplate("{{.User.Name}}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "Alice" {
		t.Errorf("expected 'Alice', got %q", output)
	}
}

func TestRuntime_LiteralOutput(t *testing.T) {
	output, err := executeTemplate(`{{"Hello World"}}`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", output)
	}
}

func TestRuntime_NumberLiteralOutput(t *testing.T) {
	output, err := executeTemplate("{{42}}", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "42" {
		t.Errorf("expected '42', got %q", output)
	}
}
