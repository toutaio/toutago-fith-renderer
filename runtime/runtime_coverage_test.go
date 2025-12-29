package runtime

import (
	"testing"

	"github.com/toutaio/toutago-fith-renderer/lexer"
	"github.com/toutaio/toutago-fith-renderer/parser"
)

// Additional tests for improved coverage

func TestRuntime_BinaryOpAllTypes(t *testing.T) {
	tests := []struct {
		expr string
		want interface{}
	}{
		// Comparison
		{"{{1 == 1}}", true},
		{"{{1 != 2}}", true},
		{"{{1 < 2}}", true},
		{"{{2 > 1}}", true},
		{"{{1 <= 1}}", true},
		{"{{2 >= 2}}", true},

		// Logical
		{"{{true && true}}", true},
		{"{{true || false}}", true},
		{"{{false && true}}", false},
		{"{{false || false}}", false},

		// Arithmetic
		{"{{5 + 3}}", 8.0},
		{"{{5 - 3}}", 2.0},
		{"{{5 * 3}}", 15.0},
		{"{{6 / 2}}", 3.0},
		{"{{7 % 3}}", 1.0},
	}

	for _, tt := range tests {
		tmpl, err := parser.New(tt.expr).Parse()
		if err != nil {
			t.Errorf("%s: parse error: %v", tt.expr, err)
			continue
		}

		rt := New()
		result, err := rt.Execute(tmpl, nil)
		if err != nil {
			t.Errorf("%s: execute error: %v", tt.expr, err)
			continue
		}

		// For numeric results, check approximately
		switch want := tt.want.(type) {
		case float64:
			if result != tt.want {
				t.Errorf("%s: expected %v, got %v", tt.expr, tt.want, result)
			}
		case bool:
			if result != want {
				t.Errorf("%s: expected %v, got %v", tt.expr, want, result)
			}
		}
	}
}

func TestRuntime_StringComparison(t *testing.T) {
	tests := []struct {
		expr string
		data map[string]interface{}
		want bool
	}{
		{`{{.a == "hello"}}`, map[string]interface{}{"a": "hello"}, true},
		{`{{.a != "world"}}`, map[string]interface{}{"a": "hello"}, true},
		{`{{.a < "zebra"}}`, map[string]interface{}{"a": "apple"}, true},
	}

	for _, tt := range tests {
		tmpl, err := parser.New(tt.expr).Parse()
		if err != nil {
			t.Errorf("%s: parse error: %v", tt.expr, err)
			continue
		}

		rt := New()
		result, err := rt.Execute(tmpl, tt.data)
		if err != nil {
			t.Errorf("%s: execute error: %v", tt.expr, err)
			continue
		}

		if result != tt.want {
			t.Errorf("%s: expected %v, got %v", tt.expr, tt.want, result)
		}
	}
}

func TestRuntime_ArithmeticErrors(t *testing.T) {
	tests := []string{
		`{{5 / 0}}`,     // Division by zero
		`{{"a" + "b"}}`, // String arithmetic
		`{{"x" * 2}}`,   // Invalid operation
	}

	for _, expr := range tests {
		tmpl, err := parser.New(expr).Parse()
		if err != nil {
			continue // Parse error is okay
		}

		rt := New()
		_, err = rt.Execute(tmpl, nil)
		if err == nil {
			t.Errorf("%s: expected error for invalid arithmetic", expr)
		}
	}
}

func TestRuntime_UnaryMinus(t *testing.T) {
	tmpl, err := parser.New("{{-5}}").Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	rt := New()
	result, err := rt.Execute(tmpl, nil)
	if err != nil {
		t.Fatalf("execute error: %v", err)
	}

	if result != -5.0 {
		t.Errorf("expected -5.0, got %v", result)
	}
}

func TestRuntime_UnaryNot(t *testing.T) {
	tests := []struct {
		expr string
		want bool
	}{
		{"{{!true}}", false},
		{"{{!false}}", true},
		{"{{!0}}", true},
		{"{{!1}}", false},
	}

	for _, tt := range tests {
		tmpl, err := parser.New(tt.expr).Parse()
		if err != nil {
			t.Errorf("%s: parse error: %v", tt.expr, err)
			continue
		}

		rt := New()
		result, err := rt.Execute(tmpl, nil)
		if err != nil {
			t.Errorf("%s: execute error: %v", tt.expr, err)
			continue
		}

		if result != tt.want {
			t.Errorf("%s: expected %v, got %v", tt.expr, tt.want, result)
		}
	}
}

func TestRuntime_IndexAccess(t *testing.T) {
	data := map[string]interface{}{
		"arr": []interface{}{10, 20, 30},
		"map": map[string]interface{}{"key": "value"},
	}

	tests := []struct {
		expr string
		want interface{}
	}{
		{"{{.arr[0]}}", 10},
		{"{{.arr[1]}}", 20},
		{"{{.map[\"key\"]}}", "value"},
	}

	for _, tt := range tests {
		tmpl, err := parser.New(tt.expr).Parse()
		if err != nil {
			t.Errorf("%s: parse error: %v", tt.expr, err)
			continue
		}

		rt := New()
		result, err := rt.Execute(tmpl, tt.data)
		if err != nil {
			t.Errorf("%s: execute error: %v", tt.expr, err)
			continue
		}

		if result != tt.want {
			t.Errorf("%s: expected %v, got %v", tt.expr, tt.want, result)
		}
	}
}

func TestRuntime_IndexOutOfBounds(t *testing.T) {
	data := map[string]interface{}{
		"arr": []interface{}{1, 2, 3},
	}

	tmpl, err := parser.New("{{.arr[10]}}").Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	rt := New()
	_, err = rt.Execute(tmpl, data)
	if err == nil {
		t.Error("expected error for out of bounds index")
	}
}

func TestRuntime_UnsupportedOperator(t *testing.T) {
	// Create a binary op node with invalid operator
	node := &parser.BinaryOpNode{
		Operator: lexer.TokenType(9999), // Invalid
		Left:     &parser.LiteralNode{Value: 1},
		Right:    &parser.LiteralNode{Value: 2},
	}

	rt := New()
	_, err := rt.evaluateBinaryOp(node)
	if err == nil {
		t.Error("expected error for unsupported operator")
	}
}

func TestRuntime_NegateInvalidValue(t *testing.T) {
	rt := New()
	_, err := rt.evaluateUnaryOp(&parser.UnaryOpNode{
		Operator: lexer.TokenMinus,
		Operand:  &parser.LiteralNode{Value: "not a number"},
	})
	if err == nil {
		t.Error("expected error when negating non-numeric value")
	}
}

func TestRuntime_TruthyValues(t *testing.T) {
	tests := []struct {
		value interface{}
		want  bool
	}{
		{true, true},
		{false, false},
		{1, true},
		{0, false},
		{0.0, false},
		{1.5, true},
		{"hello", true},
		{"", false},
		{nil, false},
		{[]interface{}{1}, true},
		{[]interface{}{}, false},
	}

	for _, tt := range tests {
		result := IsTruthy(tt.value)
		if result != tt.want {
			t.Errorf("IsTruthy(%v): expected %v, got %v", tt.value, tt.want, result)
		}
	}
}
