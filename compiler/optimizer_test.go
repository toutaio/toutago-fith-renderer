package compiler

import (
	"testing"

	"github.com/toutaio/toutago-fith-renderer/parser"
)

func TestOptimizer_OptimizeText(t *testing.T) {
	opt := NewOptimizer()
	tmpl := &parser.Template{
		Nodes: []parser.Node{
			&parser.TextNode{Value: "Hello"},
			&parser.TextNode{Value: " "},
			&parser.TextNode{Value: "World"},
		},
	}

	optimized := opt.Optimize(tmpl)

	// Text nodes should pass through unchanged
	if len(optimized.Nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(optimized.Nodes))
	}
}

func TestOptimizer_OptimizeIfTrue(t *testing.T) {
	opt := NewOptimizer()
	tmpl := &parser.Template{
		Nodes: []parser.Node{
			&parser.IfNode{
				Condition: &parser.LiteralNode{Value: true},
				Then: []parser.Node{
					&parser.TextNode{Value: "True branch"},
				},
				Else: []parser.Node{
					&parser.TextNode{Value: "False branch"},
				},
			},
		},
	}

	optimized := opt.Optimize(tmpl)

	// Should return the single then branch node directly
	if len(optimized.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(optimized.Nodes))
	}

	textNode, ok := optimized.Nodes[0].(*parser.TextNode)
	if !ok {
		t.Fatal("expected TextNode (then branch optimized)")
	}

	if textNode.Value != "True branch" {
		t.Errorf("expected 'True branch', got %q", textNode.Value)
	}
}

func TestOptimizer_OptimizeIfFalse(t *testing.T) {
	opt := NewOptimizer()
	tmpl := &parser.Template{
		Nodes: []parser.Node{
			&parser.IfNode{
				Condition: &parser.LiteralNode{Value: false},
				Then: []parser.Node{
					&parser.TextNode{Value: "True branch"},
				},
				Else: []parser.Node{
					&parser.TextNode{Value: "False branch"},
				},
			},
		},
	}

	optimized := opt.Optimize(tmpl)

	// Should return the single else branch node directly
	if len(optimized.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(optimized.Nodes))
	}

	textNode, ok := optimized.Nodes[0].(*parser.TextNode)
	if !ok {
		t.Fatal("expected TextNode (else branch optimized)")
	}

	if textNode.Value != "False branch" {
		t.Errorf("expected 'False branch', got %q", textNode.Value)
	}
}

func TestOptimizer_OptimizeIfFalseNoElse(t *testing.T) {
	opt := NewOptimizer()
	tmpl := &parser.Template{
		Nodes: []parser.Node{
			&parser.IfNode{
				Condition: &parser.LiteralNode{Value: false},
				Then: []parser.Node{
					&parser.TextNode{Value: "Never shown"},
				},
				Else: nil,
			},
		},
	}

	optimized := opt.Optimize(tmpl)

	// Entire if block should be eliminated
	if len(optimized.Nodes) != 0 {
		t.Errorf("expected 0 nodes (dead code eliminated), got %d", len(optimized.Nodes))
	}
}

func TestOptimizer_OptimizeIfDynamic(t *testing.T) {
	opt := NewOptimizer()
	tmpl := &parser.Template{
		Nodes: []parser.Node{
			&parser.IfNode{
				Condition: &parser.VariableNode{Path: []string{".", "Active"}},
				Then: []parser.Node{
					&parser.TextNode{Value: "Active"},
				},
				Else: []parser.Node{
					&parser.TextNode{Value: "Inactive"},
				},
			},
		},
	}

	optimized := opt.Optimize(tmpl)

	// Should keep if node unchanged (not constant)
	if len(optimized.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(optimized.Nodes))
	}

	ifNode, ok := optimized.Nodes[0].(*parser.IfNode)
	if !ok {
		t.Fatal("expected IfNode")
	}

	if len(ifNode.Then) != 1 {
		t.Error("then branch should have 1 node")
	}

	if len(ifNode.Else) != 1 {
		t.Error("else branch should have 1 node")
	}
}

func TestOptimizer_OptimizeRange(t *testing.T) {
	opt := NewOptimizer()
	tmpl := &parser.Template{
		Nodes: []parser.Node{
			&parser.RangeNode{
				Variable:   "item",
				Collection: &parser.VariableNode{Path: []string{".", "Items"}},
				Body: []parser.Node{
					&parser.TextNode{Value: "Item"},
				},
			},
		},
	}

	optimized := opt.Optimize(tmpl)

	// Range nodes should pass through (optimizing body)
	if len(optimized.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(optimized.Nodes))
	}

	rangeNode, ok := optimized.Nodes[0].(*parser.RangeNode)
	if !ok {
		t.Fatal("expected RangeNode")
	}

	if len(rangeNode.Body) != 1 {
		t.Error("body should have 1 node")
	}
}

func TestOptimizer_OptimizeBlock(t *testing.T) {
	opt := NewOptimizer()
	tmpl := &parser.Template{
		Nodes: []parser.Node{
			&parser.BlockNode{
				Name: "content",
				Body: []parser.Node{
					&parser.TextNode{Value: "Content"},
				},
			},
		},
	}

	optimized := opt.Optimize(tmpl)

	// Block nodes should pass through (optimizing body)
	if len(optimized.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(optimized.Nodes))
	}

	blockNode, ok := optimized.Nodes[0].(*parser.BlockNode)
	if !ok {
		t.Fatal("expected BlockNode")
	}

	if blockNode.Name != "content" {
		t.Errorf("expected name 'content', got %q", blockNode.Name)
	}

	if len(blockNode.Body) != 1 {
		t.Error("body should have 1 node")
	}
}

func TestOptimizer_OptimizeNested(t *testing.T) {
	opt := NewOptimizer()
	tmpl := &parser.Template{
		Nodes: []parser.Node{
			&parser.IfNode{
				Condition: &parser.VariableNode{Path: []string{".", "ShowContent"}},
				Then: []parser.Node{
					&parser.IfNode{
						Condition: &parser.LiteralNode{Value: true},
						Then: []parser.Node{
							&parser.TextNode{Value: "Nested content"},
						},
						Else: []parser.Node{
							&parser.TextNode{Value: "Never shown"},
						},
					},
				},
				Else: nil,
			},
		},
	}

	optimized := opt.Optimize(tmpl)

	// Outer if should remain, inner if should be optimized to just the then branch
	if len(optimized.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(optimized.Nodes))
	}

	outerIf, ok := optimized.Nodes[0].(*parser.IfNode)
	if !ok {
		t.Fatal("expected IfNode")
	}

	if len(outerIf.Then) != 1 {
		t.Fatalf("outer then should have 1 node, got %d", len(outerIf.Then))
	}

	// Inner constant true if should be optimized to just its then branch text node
	textNode, ok := outerIf.Then[0].(*parser.TextNode)
	if !ok {
		t.Fatal("expected TextNode (inner if optimized to then branch)")
	}

	if textNode.Value != "Nested content" {
		t.Errorf("expected 'Nested content', got %q", textNode.Value)
	}
}

func TestOptimizer_IsConstantBool(t *testing.T) {
	opt := NewOptimizer()

	tests := []struct {
		name     string
		node     parser.Node
		wantVal  bool
		wantBool bool
	}{
		{
			name:     "true literal",
			node:     &parser.LiteralNode{Value: true},
			wantVal:  true,
			wantBool: true,
		},
		{
			name:     "false literal",
			node:     &parser.LiteralNode{Value: false},
			wantVal:  false,
			wantBool: true,
		},
		{
			name:     "string literal",
			node:     &parser.LiteralNode{Value: "hello"},
			wantVal:  false,
			wantBool: false,
		},
		{
			name:     "variable",
			node:     &parser.VariableNode{Path: []string{".", "Active"}},
			wantVal:  false,
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, isBool := opt.isConstantBool(tt.node)
			if val != tt.wantVal {
				t.Errorf("value = %v, want %v", val, tt.wantVal)
			}
			if isBool != tt.wantBool {
				t.Errorf("isBool = %v, want %v", isBool, tt.wantBool)
			}
		})
	}
}
