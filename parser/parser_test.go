package parser

import (
	"testing"

	"github.com/toutaio/toutago-fith-renderer/lexer"
)

func TestParser_SimpleText(t *testing.T) {
	input := "Hello, World!"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ast.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(ast.Nodes))
	}

	textNode, ok := ast.Nodes[0].(*TextNode)
	if !ok {
		t.Fatalf("expected TextNode, got %T", ast.Nodes[0])
	}

	if textNode.Value != input {
		t.Errorf("expected %q, got %q", input, textNode.Value)
	}
}

func TestParser_SimpleVariable(t *testing.T) {
	input := "{{.Name}}"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ast.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(ast.Nodes))
	}

	varNode, ok := ast.Nodes[0].(*VariableNode)
	if !ok {
		t.Fatalf("expected VariableNode, got %T", ast.Nodes[0])
	}

	if len(varNode.Path) != 2 || varNode.Path[0] != "." || varNode.Path[1] != "Name" {
		t.Errorf("unexpected path: %v", varNode.Path)
	}
}

func TestParser_NestedField(t *testing.T) {
	input := "{{.User.Profile.Name}}"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	varNode, ok := ast.Nodes[0].(*VariableNode)
	if !ok {
		t.Fatalf("expected VariableNode, got %T", ast.Nodes[0])
	}

	expected := []string{".", "User", "Profile", "Name"}
	if len(varNode.Path) != len(expected) {
		t.Fatalf("expected path length %d, got %d", len(expected), len(varNode.Path))
	}

	for i, exp := range expected {
		if varNode.Path[i] != exp {
			t.Errorf("path[%d]: expected %q, got %q", i, exp, varNode.Path[i])
		}
	}
}

func TestParser_MixedTextAndVariable(t *testing.T) {
	input := "Hello {{.Name}}!"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ast.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(ast.Nodes))
	}

	// First node: "Hello "
	textNode1, ok := ast.Nodes[0].(*TextNode)
	if !ok {
		t.Fatalf("expected TextNode, got %T", ast.Nodes[0])
	}
	if textNode1.Value != "Hello " {
		t.Errorf("expected 'Hello ', got %q", textNode1.Value)
	}

	// Second node: {{.Name}}
	_, ok = ast.Nodes[1].(*VariableNode)
	if !ok {
		t.Fatalf("expected VariableNode, got %T", ast.Nodes[1])
	}

	// Third node: "!"
	textNode2, ok := ast.Nodes[2].(*TextNode)
	if !ok {
		t.Fatalf("expected TextNode, got %T", ast.Nodes[2])
	}
	if textNode2.Value != "!" {
		t.Errorf("expected '!', got %q", textNode2.Value)
	}
}

func TestParser_StringLiteral(t *testing.T) {
	input := `{{"hello"}}`
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	litNode, ok := ast.Nodes[0].(*LiteralNode)
	if !ok {
		t.Fatalf("expected LiteralNode, got %T", ast.Nodes[0])
	}

	if litNode.Value != "hello" {
		t.Errorf("expected 'hello', got %v", litNode.Value)
	}
}

func TestParser_NumberLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"{{42}}", 42},
		{"{{123}}", 123},
		{"{{3.14}}", 3.14},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			ast, err := p.Parse()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			litNode, ok := ast.Nodes[0].(*LiteralNode)
			if !ok {
				t.Fatalf("expected LiteralNode, got %T", ast.Nodes[0])
			}

			if litNode.Value != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, litNode.Value)
			}
		})
	}
}

func TestParser_FunctionCall(t *testing.T) {
	input := "{{upper .Name}}"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	callNode, ok := ast.Nodes[0].(*CallNode)
	if !ok {
		t.Fatalf("expected CallNode, got %T", ast.Nodes[0])
	}

	if callNode.Function != "upper" {
		t.Errorf("expected function 'upper', got %q", callNode.Function)
	}

	if len(callNode.Args) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(callNode.Args))
	}
}

func TestParser_FunctionWithMultipleArgs(t *testing.T) {
	input := "{{truncate .Text 100}}"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	callNode, ok := ast.Nodes[0].(*CallNode)
	if !ok {
		t.Fatalf("expected CallNode, got %T", ast.Nodes[0])
	}

	if callNode.Function != "truncate" {
		t.Errorf("expected function 'truncate', got %q", callNode.Function)
	}

	if len(callNode.Args) != 2 {
		t.Fatalf("expected 2 arguments, got %d", len(callNode.Args))
	}
}

func TestParser_Pipe(t *testing.T) {
	input := "{{.Name | upper | trim}}"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pipeNode, ok := ast.Nodes[0].(*PipeNode)
	if !ok {
		t.Fatalf("expected PipeNode, got %T", ast.Nodes[0])
	}

	// Value should be a VariableNode
	_, ok = pipeNode.Value.(*VariableNode)
	if !ok {
		t.Fatalf("expected VariableNode as value, got %T", pipeNode.Value)
	}

	// Should have two filters
	if len(pipeNode.Filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(pipeNode.Filters))
	}

	if pipeNode.Filters[0] != "upper" {
		t.Errorf("expected filter 'upper', got %q", pipeNode.Filters[0])
	}

	if pipeNode.Filters[1] != "trim" {
		t.Errorf("expected filter 'trim', got %q", pipeNode.Filters[1])
	}
}

func TestParser_BinaryOp(t *testing.T) {
	input := "{{.A + .B}}"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	binOp, ok := ast.Nodes[0].(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode, got %T", ast.Nodes[0])
	}

	if binOp.Operator != lexer.TokenPlus {
		t.Errorf("expected + operator, got %v", binOp.Operator)
	}

	// Left should be .A
	leftVar, ok := binOp.Left.(*VariableNode)
	if !ok {
		t.Fatalf("expected VariableNode on left, got %T", binOp.Left)
	}
	if leftVar.Path[1] != "A" {
		t.Errorf("expected .A on left, got %v", leftVar.Path)
	}

	// Right should be .B
	rightVar, ok := binOp.Right.(*VariableNode)
	if !ok {
		t.Fatalf("expected VariableNode on right, got %T", binOp.Right)
	}
	if rightVar.Path[1] != "B" {
		t.Errorf("expected .B on right, got %v", rightVar.Path)
	}
}

func TestParser_ComparisonOp(t *testing.T) {
	tests := []struct {
		input    string
		operator lexer.TokenType
	}{
		{"{{.A == .B}}", lexer.TokenEqual},
		{"{{.A != .B}}", lexer.TokenNotEqual},
		{"{{.A < .B}}", lexer.TokenLess},
		{"{{.A > .B}}", lexer.TokenGreater},
		{"{{.A <= .B}}", lexer.TokenLessEq},
		{"{{.A >= .B}}", lexer.TokenGreaterEq},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			ast, err := p.Parse()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			binOp, ok := ast.Nodes[0].(*BinaryOpNode)
			if !ok {
				t.Fatalf("expected BinaryOpNode, got %T", ast.Nodes[0])
			}

			if binOp.Operator != tt.operator {
				t.Errorf("expected operator %v, got %v", tt.operator, binOp.Operator)
			}
		})
	}
}

func TestParser_UnaryOp(t *testing.T) {
	input := "{{!.Active}}"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	unaryOp, ok := ast.Nodes[0].(*UnaryOpNode)
	if !ok {
		t.Fatalf("expected UnaryOpNode, got %T", ast.Nodes[0])
	}

	if unaryOp.Operator != lexer.TokenNot {
		t.Errorf("expected ! operator, got %v", unaryOp.Operator)
	}
}

func TestParser_If(t *testing.T) {
	input := "{{if .Active}}yes{{end}}"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ifNode, ok := ast.Nodes[0].(*IfNode)
	if !ok {
		t.Fatalf("expected IfNode, got %T", ast.Nodes[0])
	}

	// Condition should be a VariableNode
	_, ok = ifNode.Condition.(*VariableNode)
	if !ok {
		t.Fatalf("expected VariableNode as condition, got %T", ifNode.Condition)
	}

	// Then body should have one text node
	if len(ifNode.Then) != 1 {
		t.Fatalf("expected 1 node in then body, got %d", len(ifNode.Then))
	}

	textNode, ok := ifNode.Then[0].(*TextNode)
	if !ok {
		t.Fatalf("expected TextNode in then body, got %T", ifNode.Then[0])
	}

	if textNode.Value != "yes" {
		t.Errorf("expected 'yes', got %q", textNode.Value)
	}

	// Else body should be nil
	if ifNode.Else != nil {
		t.Errorf("expected nil else body, got %v", ifNode.Else)
	}
}

func TestParser_IfElse(t *testing.T) {
	input := "{{if .Active}}yes{{else}}no{{end}}"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ifNode, ok := ast.Nodes[0].(*IfNode)
	if !ok {
		t.Fatalf("expected IfNode, got %T", ast.Nodes[0])
	}

	// Then body
	if len(ifNode.Then) != 1 {
		t.Fatalf("expected 1 node in then body, got %d", len(ifNode.Then))
	}

	textNode1, ok := ifNode.Then[0].(*TextNode)
	if !ok || textNode1.Value != "yes" {
		t.Errorf("expected 'yes' in then body")
	}

	// Else body
	if len(ifNode.Else) != 1 {
		t.Fatalf("expected 1 node in else body, got %d", len(ifNode.Else))
	}

	textNode2, ok := ifNode.Else[0].(*TextNode)
	if !ok || textNode2.Value != "no" {
		t.Errorf("expected 'no' in else body")
	}
}

func TestParser_Range(t *testing.T) {
	input := "{{range .Items}}{{.}}{{end}}"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rangeNode, ok := ast.Nodes[0].(*RangeNode)
	if !ok {
		t.Fatalf("expected RangeNode, got %T", ast.Nodes[0])
	}

	// Collection should be a VariableNode
	_, ok = rangeNode.Collection.(*VariableNode)
	if !ok {
		t.Fatalf("expected VariableNode as collection, got %T", rangeNode.Collection)
	}

	// Body should have one variable node
	if len(rangeNode.Body) != 1 {
		t.Fatalf("expected 1 node in body, got %d", len(rangeNode.Body))
	}
}

func TestParser_Include(t *testing.T) {
	input := `{{include "header"}}`
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	includeNode, ok := ast.Nodes[0].(*IncludeNode)
	if !ok {
		t.Fatalf("expected IncludeNode, got %T", ast.Nodes[0])
	}

	if includeNode.Template != "header" {
		t.Errorf("expected template 'header', got %q", includeNode.Template)
	}
}

func TestParser_Extends(t *testing.T) {
	input := `{{extends "layout"}}`
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	extendsNode, ok := ast.Nodes[0].(*ExtendsNode)
	if !ok {
		t.Fatalf("expected ExtendsNode, got %T", ast.Nodes[0])
	}

	if extendsNode.Template != "layout" {
		t.Errorf("expected template 'layout', got %q", extendsNode.Template)
	}
}

func TestParser_Block(t *testing.T) {
	input := `{{block "content"}}default{{end}}`
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	blockNode, ok := ast.Nodes[0].(*BlockNode)
	if !ok {
		t.Fatalf("expected BlockNode, got %T", ast.Nodes[0])
	}

	if blockNode.Name != "content" {
		t.Errorf("expected block name 'content', got %q", blockNode.Name)
	}

	if len(blockNode.Body) != 1 {
		t.Fatalf("expected 1 node in body, got %d", len(blockNode.Body))
	}
}

func TestParser_ArrayAccess(t *testing.T) {
	input := "{{.Items[0]}}"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	indexNode, ok := ast.Nodes[0].(*IndexNode)
	if !ok {
		t.Fatalf("expected IndexNode, got %T", ast.Nodes[0])
	}

	// Object should be VariableNode for .Items
	varNode, ok := indexNode.Object.(*VariableNode)
	if !ok {
		t.Fatalf("expected VariableNode as object, got %T", indexNode.Object)
	}

	if varNode.Path[1] != "Items" {
		t.Errorf("expected .Items, got %v", varNode.Path)
	}

	// Index should be a literal 0
	litNode, ok := indexNode.Index.(*LiteralNode)
	if !ok {
		t.Fatalf("expected LiteralNode as index, got %T", indexNode.Index)
	}

	if litNode.Value != 0 {
		t.Errorf("expected index 0, got %v", litNode.Value)
	}
}

func TestParser_MapAccess(t *testing.T) {
	input := `{{.Data["key"]}}`
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	indexNode, ok := ast.Nodes[0].(*IndexNode)
	if !ok {
		t.Fatalf("expected IndexNode, got %T", ast.Nodes[0])
	}

	// Index should be string literal "key"
	litNode, ok := indexNode.Index.(*LiteralNode)
	if !ok {
		t.Fatalf("expected LiteralNode as index, got %T", indexNode.Index)
	}

	if litNode.Value != "key" {
		t.Errorf("expected index 'key', got %v", litNode.Value)
	}
}

func TestParser_EmptyTemplate(t *testing.T) {
	input := ""
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ast.Nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(ast.Nodes))
	}
}

func TestParser_ComplexTemplate(t *testing.T) {
	input := `
<h1>{{.Title}}</h1>
{{if .User.Active}}
	<p>Welcome, {{.User.Name | upper}}!</p>
	{{range .Items}}
		<li>{{.}}</li>
	{{end}}
{{end}}
`
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have multiple nodes
	if len(ast.Nodes) == 0 {
		t.Error("expected nodes in complex template")
	}

	// Find the if node
	var foundIf bool
	for _, node := range ast.Nodes {
		if _, ok := node.(*IfNode); ok {
			foundIf = true
			break
		}
	}

	if !foundIf {
		t.Error("expected to find IfNode in complex template")
	}
}

func TestParser_GroupedExpression(t *testing.T) {
	input := "{{(.A + .B) * .C}}"
	l := lexer.New(input)
	p := New(l)

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should parse as multiplication with grouped addition on left
	binOp, ok := ast.Nodes[0].(*BinaryOpNode)
	if !ok {
		t.Fatalf("expected BinaryOpNode, got %T", ast.Nodes[0])
	}

	if binOp.Operator != lexer.TokenMult {
		t.Errorf("expected * operator, got %v", binOp.Operator)
	}
}
