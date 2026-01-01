package parser

import "testing"

// TestASTNodeMethods tests the Pos() and String() methods of all AST nodes
func TestASTNodeMethods(t *testing.T) {
	tests := []struct {
		name string
		node Node
	}{
		{"Template", &Template{Nodes: []Node{}}},
		{"TextNode", &TextNode{Value: "test", Position: Position{Line: 1, Column: 1}}},
		{
			"VariableNode",
			&VariableNode{Path: []string{"test"}, Position: Position{Line: 1, Column: 1}},
		},
		{
			"BinaryOpNode",
			&BinaryOpNode{
				Left:     &LiteralNode{Value: 1},
				Right:    &LiteralNode{Value: 2},
				Position: Position{Line: 1, Column: 1},
			},
		},
		{
			"UnaryOpNode",
			&UnaryOpNode{Operand: &LiteralNode{Value: true}, Position: Position{Line: 1, Column: 1}},
		},
		{"LiteralNode", &LiteralNode{Value: "test", Position: Position{Line: 1, Column: 1}}},
		{
			"IndexNode",
			&IndexNode{
				Object:   &VariableNode{Path: []string{"items"}},
				Index:    &LiteralNode{Value: 0},
				Position: Position{Line: 1, Column: 1},
			},
		},
		{"CallNode", &CallNode{Function: "upper", Args: []Node{}, Position: Position{Line: 1, Column: 1}}},
		{
			"PipeNode",
			&PipeNode{
				Value:    &VariableNode{Path: []string{"test"}},
				Filters:  []string{"upper"},
				Position: Position{Line: 1, Column: 1},
			},
		},
		{
			"IfNode",
			&IfNode{
				Condition: &LiteralNode{Value: true},
				Then:      []Node{},
				Else:      []Node{},
				Position:  Position{Line: 1, Column: 1},
			},
		},
		{
			"RangeNode",
			&RangeNode{
				Variable:   "item",
				Collection: &VariableNode{Path: []string{"items"}},
				Body:       []Node{},
				Position:   Position{Line: 1, Column: 1},
			},
		},
		{"IncludeNode", &IncludeNode{Template: "partial", Position: Position{Line: 1, Column: 1}}},
		{"ExtendsNode", &ExtendsNode{Template: "base", Position: Position{Line: 1, Column: 1}}},
		{"BlockNode", &BlockNode{Name: "content", Body: []Node{}, Position: Position{Line: 1, Column: 1}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := tt.node.Pos()
			if pos.Line < 0 || pos.Column < 0 {
				t.Errorf("%s.Pos() returned negative values: %+v", tt.name, pos)
			}
			str := tt.node.String()
			if str == "" {
				t.Errorf("%s.String() returned empty string", tt.name)
			}
		})
	}
}
