package parser

import "github.com/toutaio/toutago-fith-renderer/lexer"

// Node represents a node in the Abstract Syntax Tree.
// All AST node types implement this interface.
type Node interface {
	// Pos returns the position of the node in the source template.
	Pos() Position
	// String returns a string representation for debugging.
	String() string
}

// Position represents a position in the source template.
type Position struct {
	Line   int // Line number (1-indexed)
	Column int // Column number (1-indexed)
}

// Template is the root node of the AST.
// It contains a list of top-level nodes (text, expressions, statements).
type Template struct {
	Nodes []Node
}

func (t *Template) Pos() Position { return Position{Line: 1, Column: 1} }
func (t *Template) String() string {
	return "Template"
}

// TextNode represents literal text content.
type TextNode struct {
	Position Position
	Value    string
}

func (n *TextNode) Pos() Position  { return n.Position }
func (n *TextNode) String() string { return "Text: " + n.Value }

// VariableNode represents a variable expression like {{.Name}} or {{.User.Email}}.
type VariableNode struct {
	Position Position
	Path     []string // Path components, e.g., [".", "User", "Email"]
}

func (n *VariableNode) Pos() Position  { return n.Position }
func (n *VariableNode) String() string { return "Variable" }

// BinaryOpNode represents a binary operation like {{.A + .B}}.
type BinaryOpNode struct {
	Position Position
	Operator lexer.TokenType // The operator (+, -, *, /, etc.)
	Left     Node
	Right    Node
}

func (n *BinaryOpNode) Pos() Position  { return n.Position }
func (n *BinaryOpNode) String() string { return "BinaryOp" }

// UnaryOpNode represents a unary operation like {{!.Active}}.
type UnaryOpNode struct {
	Position Position
	Operator lexer.TokenType // The operator (!, -)
	Operand  Node
}

func (n *UnaryOpNode) Pos() Position  { return n.Position }
func (n *UnaryOpNode) String() string { return "UnaryOp" }

// LiteralNode represents a literal value (string, number, boolean).
type LiteralNode struct {
	Position Position
	Value    interface{} // The actual value (string, int, float64, bool)
}

func (n *LiteralNode) Pos() Position  { return n.Position }
func (n *LiteralNode) String() string { return "Literal" }

// IndexNode represents array/map access like {{.Items[0]}} or {{.Data["key"]}}.
type IndexNode struct {
	Position Position
	Object   Node // The object being indexed
	Index    Node // The index expression
}

func (n *IndexNode) Pos() Position  { return n.Position }
func (n *IndexNode) String() string { return "Index" }

// CallNode represents a function call like {{upper .Name}}.
type CallNode struct {
	Position Position
	Function string // Function name
	Args     []Node // Arguments
}

func (n *CallNode) Pos() Position  { return n.Position }
func (n *CallNode) String() string { return "Call: " + n.Function }

// PipeNode represents a filter pipeline like {{.Name | upper | trim}}.
type PipeNode struct {
	Position Position
	Value    Node     // Initial value
	Filters  []string // List of filter function names
}

func (n *PipeNode) Pos() Position  { return n.Position }
func (n *PipeNode) String() string { return "Pipe" }

// IfNode represents an if/else statement.
type IfNode struct {
	Position  Position
	Condition Node   // The condition expression
	Then      []Node // Nodes to execute if condition is true
	Else      []Node // Nodes to execute if condition is false (may be nil)
}

func (n *IfNode) Pos() Position  { return n.Position }
func (n *IfNode) String() string { return "If" }

// RangeNode represents a range loop.
type RangeNode struct {
	Position   Position
	Variable   string // Loop variable name (e.g., "item")
	KeyVar     string // Key variable for map iteration (optional)
	Collection Node   // The collection to iterate over
	Body       []Node // Nodes to execute for each iteration
}

func (n *RangeNode) Pos() Position  { return n.Position }
func (n *RangeNode) String() string { return "Range" }

// IncludeNode represents an include directive.
type IncludeNode struct {
	Position Position
	Template string          // Template name to include
	Params   map[string]Node // Parameters to pass
	Context  Node            // Context to pass (may be nil)
}

func (n *IncludeNode) Pos() Position  { return n.Position }
func (n *IncludeNode) String() string { return "Include: " + n.Template }

// ExtendsNode represents a layout extension.
type ExtendsNode struct {
	Position Position
	Template string // Parent template name
}

func (n *ExtendsNode) Pos() Position  { return n.Position }
func (n *ExtendsNode) String() string { return "Extends: " + n.Template }

// BlockNode represents a named block that can be overridden.
type BlockNode struct {
	Position Position
	Name     string // Block name
	Body     []Node // Default content
}

func (n *BlockNode) Pos() Position  { return n.Position }
func (n *BlockNode) String() string { return "Block: " + n.Name }
