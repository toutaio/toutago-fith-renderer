package compiler

import (
	"github.com/toutaio/toutago-fith-renderer/parser"
)

// Optimizer performs AST optimizations.
type Optimizer struct{}

// NewOptimizer creates a new optimizer.
func NewOptimizer() *Optimizer {
	return &Optimizer{}
}

// Optimize performs all optimization passes on the template AST.
func (o *Optimizer) Optimize(tmpl *parser.Template) *parser.Template {
	optimized := &parser.Template{
		Nodes: make([]parser.Node, 0, len(tmpl.Nodes)),
	}

	for _, node := range tmpl.Nodes {
		if optimizedNode := o.optimizeNode(node); optimizedNode != nil {
			optimized.Nodes = append(optimized.Nodes, optimizedNode)
		}
	}

	return optimized
}

// optimizeNode optimizes a single node.
func (o *Optimizer) optimizeNode(node parser.Node) parser.Node {
	switch n := node.(type) {
	case *parser.TextNode:
		return n

	case *parser.VariableNode:
		return n

	case *parser.IfNode:
		return o.optimizeIf(n)

	case *parser.RangeNode:
		return o.optimizeRange(n)

	case *parser.IncludeNode:
		return n

	case *parser.ExtendsNode:
		return n

	case *parser.BlockNode:
		return o.optimizeBlock(n)

	default:
		return node
	}
}

// optimizeIf optimizes if nodes, performing constant folding and dead code elimination.
func (o *Optimizer) optimizeIf(n *parser.IfNode) parser.Node {
	// Check if condition is a constant boolean
	if constVal, isConst := o.isConstantBool(n.Condition); isConst {
		if constVal {
			// Condition is always true - return then branch
			if len(n.Then) == 1 {
				return o.optimizeNode(n.Then[0])
			}
			// Multiple nodes - keep if structure but mark as optimized
			return &parser.IfNode{
				Position:  n.Position,
				Condition: n.Condition,
				Then:      o.optimizeNodes(n.Then),
				Else:      nil, // Dead code eliminated
			}
		}
		// Condition is always false - return else branch or nothing
		if len(n.Else) == 1 {
			return o.optimizeNode(n.Else[0])
		}
		if len(n.Else) > 0 {
			return &parser.IfNode{
				Position:  n.Position,
				Condition: &parser.LiteralNode{Position: n.Position, Value: false},
				Then:      nil, // Dead code eliminated
				Else:      o.optimizeNodes(n.Else),
			}
		}
		return nil // Entire if block eliminated
	}

	// Not constant - optimize branches
	var elseNodes []parser.Node
	if n.Else != nil {
		elseNodes = o.optimizeNodes(n.Else)
	}
	return &parser.IfNode{
		Position:  n.Position,
		Condition: n.Condition,
		Then:      o.optimizeNodes(n.Then),
		Else:      elseNodes,
	}
}

// optimizeRange optimizes range nodes.
func (o *Optimizer) optimizeRange(n *parser.RangeNode) parser.Node {
	// Optimize loop body
	return &parser.RangeNode{
		Position:   n.Position,
		Variable:   n.Variable,
		KeyVar:     n.KeyVar,
		Collection: n.Collection,
		Body:       o.optimizeNodes(n.Body),
	}
}

// optimizeBlock optimizes block nodes.
func (o *Optimizer) optimizeBlock(n *parser.BlockNode) parser.Node {
	return &parser.BlockNode{
		Position: n.Position,
		Name:     n.Name,
		Body:     o.optimizeNodes(n.Body),
	}
}

// optimizeNodes optimizes a slice of nodes.
func (o *Optimizer) optimizeNodes(nodes []parser.Node) []parser.Node {
	optimized := make([]parser.Node, 0, len(nodes))
	for _, node := range nodes {
		if opt := o.optimizeNode(node); opt != nil {
			optimized = append(optimized, opt)
		}
	}
	return optimized
}

// isConstantBool checks if a node is a constant boolean literal.
// Returns (value, true) if it's a constant bool, or (false, false) otherwise.
func (o *Optimizer) isConstantBool(node parser.Node) (value, ok bool) {
	if lit, ok := node.(*parser.LiteralNode); ok {
		if b, ok := lit.Value.(bool); ok {
			return b, true
		}
	}
	return false, false
}
