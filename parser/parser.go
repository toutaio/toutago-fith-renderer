package parser

import (
	"fmt"
	"strconv"

	"github.com/toutaio/toutago-fith-renderer/lexer"
)

// Parser builds an Abstract Syntax Tree from tokens.
type Parser struct {
	lexer   *lexer.Lexer
	current lexer.Token
	peek    lexer.Token
}

// New creates a new Parser for the given lexer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	// Read two tokens so current and peek are both set
	p.nextToken()
	p.nextToken()
	return p
}

// Parse parses the entire template and returns the AST root.
func (p *Parser) Parse() (*Template, error) {
	template := &Template{
		Nodes: []Node{},
	}

	for p.current.Type != lexer.TokenEOF {
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if node != nil {
			template.Nodes = append(template.Nodes, node)
		}
	}

	return template, nil
}

// parseNode parses a single top-level node.
func (p *Parser) parseNode() (Node, error) {
	switch p.current.Type {
	case lexer.TokenText:
		return p.parseText()
	case lexer.TokenOpenDelim:
		return p.parseExpression()
	default:
		return nil, p.error(fmt.Sprintf("unexpected token: %v", p.current.Type))
	}
}

// parseText parses a text node.
func (p *Parser) parseText() (Node, error) {
	node := &TextNode{
		Position: Position{Line: p.current.Line, Column: p.current.Column},
		Value:    p.current.Value,
	}
	p.nextToken()
	return node, nil
}

// parseExpression parses an expression starting with {{.
func (p *Parser) parseExpression() (Node, error) {
	if p.current.Type != lexer.TokenOpenDelim {
		return nil, p.error("expected {{")
	}
	p.nextToken() // consume {{

	// Check for keywords first
	switch p.current.Type {
	case lexer.TokenIf:
		return p.parseIf()
	case lexer.TokenRange:
		return p.parseRange()
	case lexer.TokenInclude:
		return p.parseInclude()
	case lexer.TokenExtends:
		return p.parseExtends()
	case lexer.TokenBlock:
		return p.parseBlock()
	case lexer.TokenEnd:
		// End token without matching start - this is an error
		return nil, p.error("unexpected 'end' token")
	case lexer.TokenElse:
		// Else token outside if - this is an error
		return nil, p.error("unexpected 'else' token")
	}

	// Otherwise, parse as a value expression
	node, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	// Expect closing delimiter
	if p.current.Type != lexer.TokenCloseDelim {
		return nil, p.error(fmt.Sprintf("expected }}, got %v", p.current.Type))
	}
	p.nextToken() // consume }}

	return node, nil
}

// parseValue parses a value expression (variable, literal, function call, etc.).
func (p *Parser) parseValue() (Node, error) {
	// Start with primary expression
	node, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	// Check for pipe operator
	if p.current.Type == lexer.TokenPipe {
		return p.parsePipe(node)
	}

	// Check for binary operators
	if p.isBinaryOp(p.current.Type) {
		return p.parseBinaryOp(node)
	}

	return node, nil
}

// parsePrimary parses a primary expression (variable, literal, function call, etc.).
func (p *Parser) parsePrimary() (Node, error) {
	pos := Position{Line: p.current.Line, Column: p.current.Column}

	switch p.current.Type {
	case lexer.TokenDot:
		return p.parseVariable()
	case lexer.TokenIdent:
		// Could be a function call
		return p.parseFunctionCall()
	case lexer.TokenString:
		value := p.current.Value
		p.nextToken()
		return &LiteralNode{Position: pos, Value: value}, nil
	case lexer.TokenNumber:
		return p.parseNumber(pos)
	case lexer.TokenNot:
		return p.parseUnaryOp()
	case lexer.TokenLParen:
		return p.parseGrouped()
	default:
		return nil, p.error(fmt.Sprintf("unexpected token in expression: %v", p.current.Type))
	}
}

// parseVariable parses a variable expression like .Name or .User.Email.
func (p *Parser) parseVariable() (Node, error) {
	pos := Position{Line: p.current.Line, Column: p.current.Column}
	path := []string{"."}

	p.nextToken() // consume initial dot

	// Parse field access chain
	for p.current.Type == lexer.TokenIdent {
		path = append(path, p.current.Value)
		p.nextToken()

		// Check for array/map access
		if p.current.Type == lexer.TokenLBrack {
			// Convert to IndexNode
			varNode := &VariableNode{Position: pos, Path: path}
			return p.parseIndex(varNode)
		}

		// Check for continued dot notation
		if p.current.Type == lexer.TokenDot {
			p.nextToken()
		} else {
			break
		}
	}

	return &VariableNode{Position: pos, Path: path}, nil
}

// parseIndex parses array/map access like [0] or ["key"].
func (p *Parser) parseIndex(object Node) (Node, error) {
	pos := Position{Line: p.current.Line, Column: p.current.Column}

	p.nextToken() // consume [

	index, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	if p.current.Type != lexer.TokenRBrack {
		return nil, p.error("expected ]")
	}
	p.nextToken() // consume ]

	return &IndexNode{Position: pos, Object: object, Index: index}, nil
}

// parseFunctionCall parses a function call like upper .Name or truncate .Text 100.
func (p *Parser) parseFunctionCall() (Node, error) {
	pos := Position{Line: p.current.Line, Column: p.current.Column}
	funcName := p.current.Value
	p.nextToken()

	args := []Node{}

	// Parse arguments until we hit }} or |
	for p.current.Type != lexer.TokenCloseDelim && p.current.Type != lexer.TokenPipe {
		arg, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		// Skip optional commas between arguments
		if p.current.Type == lexer.TokenComma {
			p.nextToken()
		}
	}

	return &CallNode{Position: pos, Function: funcName, Args: args}, nil
}

// parsePipe parses a pipe expression like .Name | upper | trim.
func (p *Parser) parsePipe(value Node) (Node, error) {
	pos := Position{Line: p.current.Line, Column: p.current.Column}
	filters := []string{}

	for p.current.Type == lexer.TokenPipe {
		p.nextToken() // consume |

		if p.current.Type != lexer.TokenIdent {
			return nil, p.error("expected filter name after |")
		}

		filters = append(filters, p.current.Value)
		p.nextToken()
	}

	return &PipeNode{Position: pos, Value: value, Filters: filters}, nil
}

// parseBinaryOp parses a binary operation.
func (p *Parser) parseBinaryOp(left Node) (Node, error) {
	pos := Position{Line: p.current.Line, Column: p.current.Column}
	op := p.current.Type
	p.nextToken()

	right, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	return &BinaryOpNode{Position: pos, Operator: op, Left: left, Right: right}, nil
}

// parseUnaryOp parses a unary operation like !.Active.
func (p *Parser) parseUnaryOp() (Node, error) {
	pos := Position{Line: p.current.Line, Column: p.current.Column}
	op := p.current.Type
	p.nextToken()

	operand, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	return &UnaryOpNode{Position: pos, Operator: op, Operand: operand}, nil
}

// parseGrouped parses a grouped expression like (.A + .B).
func (p *Parser) parseGrouped() (Node, error) {
	p.nextToken() // consume (

	node, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	if p.current.Type != lexer.TokenRParen {
		return nil, p.error("expected )")
	}
	p.nextToken() // consume )

	return node, nil
}

// parseNumber parses a number literal.
func (p *Parser) parseNumber(pos Position) (Node, error) {
	value := p.current.Value
	p.nextToken()

	// Try to parse as int first
	if intVal, err := strconv.Atoi(value); err == nil {
		return &LiteralNode{Position: pos, Value: intVal}, nil
	}

	// Parse as float
	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, p.error(fmt.Sprintf("invalid number: %s", value))
	}

	return &LiteralNode{Position: pos, Value: floatVal}, nil
}

// parseIf parses an if statement.
func (p *Parser) parseIf() (Node, error) {
	pos := Position{Line: p.current.Line, Column: p.current.Column}
	p.nextToken() // consume 'if'

	// Parse condition
	condition, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	// Expect }}
	if p.current.Type != lexer.TokenCloseDelim {
		return nil, p.error("expected }} after if condition")
	}
	p.nextToken() // consume }}

	// Parse then body
	thenBody, err := p.parseUntil(lexer.TokenElse, lexer.TokenEnd)
	if err != nil {
		return nil, err
	}

	var elseBody []Node

	// Check for else clause
	if p.current.Type == lexer.TokenOpenDelim && p.peek.Type == lexer.TokenElse {
		p.nextToken() // consume {{
		p.nextToken() // consume 'else'

		if p.current.Type != lexer.TokenCloseDelim {
			return nil, p.error("expected }} after else")
		}
		p.nextToken() // consume }}

		// Parse else body
		elseBody, err = p.parseUntil(lexer.TokenEnd)
		if err != nil {
			return nil, err
		}
	}

	// Expect {{end}}
	if p.current.Type != lexer.TokenOpenDelim || p.peek.Type != lexer.TokenEnd {
		return nil, p.error("expected {{end}}")
	}
	p.nextToken() // consume {{
	p.nextToken() // consume 'end'

	if p.current.Type != lexer.TokenCloseDelim {
		return nil, p.error("expected }} after end")
	}
	p.nextToken() // consume }}

	return &IfNode{Position: pos, Condition: condition, Then: thenBody, Else: elseBody}, nil
}

// parseRange parses a range loop (simplified for now).
func (p *Parser) parseRange() (Node, error) {
	pos := Position{Line: p.current.Line, Column: p.current.Column}
	p.nextToken() // consume 'range'

	// Parse collection
	collection, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	// Expect }}
	if p.current.Type != lexer.TokenCloseDelim {
		return nil, p.error("expected }} after range")
	}
	p.nextToken() // consume }}

	// Parse body
	body, err := p.parseUntil(lexer.TokenEnd)
	if err != nil {
		return nil, err
	}

	// Expect {{end}}
	if p.current.Type != lexer.TokenOpenDelim || p.peek.Type != lexer.TokenEnd {
		return nil, p.error("expected {{end}}")
	}
	p.nextToken() // consume {{
	p.nextToken() // consume 'end'

	if p.current.Type != lexer.TokenCloseDelim {
		return nil, p.error("expected }} after end")
	}
	p.nextToken() // consume }}

	return &RangeNode{Position: pos, Collection: collection, Body: body}, nil
}

// parseInclude parses an include directive with optional parameters.
// Supports: {{include "template"}} or {{include "template" key=value}} or {{include "template" .context}}
func (p *Parser) parseInclude() (Node, error) {
	pos := Position{Line: p.current.Line, Column: p.current.Column}
	p.nextToken() // consume 'include'

	if p.current.Type != lexer.TokenString {
		return nil, p.error("expected template name after include")
	}

	templateName := p.current.Value
	p.nextToken()

	// Parse optional parameters or context
	var params map[string]Node
	var context Node

	// Check if there are more tokens before }}
	for p.current.Type != lexer.TokenCloseDelim && p.current.Type != lexer.TokenEOF {
		// Check if it's a named parameter (key=value)
		if p.current.Type == lexer.TokenIdent && p.peek.Type == lexer.TokenAssign {
			if params == nil {
				params = make(map[string]Node)
			}

			paramName := p.current.Value
			p.nextToken() // consume param name
			p.nextToken() // consume '='

			// Parse the value expression (use parsePrimary for simple values)
			valueExpr, err := p.parsePrimary()
			if err != nil {
				return nil, err
			}
			params[paramName] = valueExpr
		} else {
			// It's a context expression (like .user)
			ctxExpr, err := p.parseValue()
			if err != nil {
				return nil, err
			}
			context = ctxExpr
			break // Context should be last
		}
	}

	// Expect }}
	if p.current.Type != lexer.TokenCloseDelim {
		return nil, p.error("expected }} after include")
	}
	p.nextToken() // consume }}

	return &IncludeNode{
		Position: pos,
		Template: templateName,
		Params:   params,
		Context:  context,
	}, nil
}

// parseExtends parses an extends directive.
func (p *Parser) parseExtends() (Node, error) {
	pos := Position{Line: p.current.Line, Column: p.current.Column}
	p.nextToken() // consume 'extends'

	if p.current.Type != lexer.TokenString {
		return nil, p.error("expected template name after extends")
	}

	templateName := p.current.Value
	p.nextToken()

	// Expect }}
	if p.current.Type != lexer.TokenCloseDelim {
		return nil, p.error("expected }} after extends")
	}
	p.nextToken() // consume }}

	return &ExtendsNode{Position: pos, Template: templateName}, nil
}

// parseBlock parses a block directive.
func (p *Parser) parseBlock() (Node, error) {
	pos := Position{Line: p.current.Line, Column: p.current.Column}
	p.nextToken() // consume 'block'

	if p.current.Type != lexer.TokenString {
		return nil, p.error("expected block name")
	}

	blockName := p.current.Value
	p.nextToken()

	// Expect }}
	if p.current.Type != lexer.TokenCloseDelim {
		return nil, p.error("expected }} after block name")
	}
	p.nextToken() // consume }}

	// Parse body
	body, err := p.parseUntil(lexer.TokenEnd)
	if err != nil {
		return nil, err
	}

	// Expect {{end}}
	if p.current.Type != lexer.TokenOpenDelim || p.peek.Type != lexer.TokenEnd {
		return nil, p.error("expected {{end}}")
	}
	p.nextToken() // consume {{
	p.nextToken() // consume 'end'

	if p.current.Type != lexer.TokenCloseDelim {
		return nil, p.error("expected }} after end")
	}
	p.nextToken() // consume }}

	return &BlockNode{Position: pos, Name: blockName, Body: body}, nil
}

// parseUntil parses nodes until one of the specified token types is encountered.
func (p *Parser) parseUntil(stopTokens ...lexer.TokenType) ([]Node, error) {
	nodes := []Node{}

	for p.current.Type != lexer.TokenEOF {
		// Check if we've hit a stop token (which would be after {{)
		if p.current.Type == lexer.TokenOpenDelim {
			for _, stop := range stopTokens {
				if p.peek.Type == stop {
					return nodes, nil
				}
			}
		}

		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if node != nil {
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

// Helper methods

func (p *Parser) nextToken() {
	p.current = p.peek
	tok, err := p.lexer.NextToken()
	if err != nil {
		// Store error in peek token
		p.peek = lexer.Token{Type: lexer.TokenError, Value: err.Error()}
	} else {
		p.peek = tok
	}
}

func (p *Parser) isBinaryOp(t lexer.TokenType) bool {
	return t == lexer.TokenPlus || t == lexer.TokenMinus ||
		t == lexer.TokenMult || t == lexer.TokenDiv || t == lexer.TokenMod ||
		t == lexer.TokenEqual || t == lexer.TokenNotEqual ||
		t == lexer.TokenLess || t == lexer.TokenGreater ||
		t == lexer.TokenLessEq || t == lexer.TokenGreaterEq ||
		t == lexer.TokenAnd || t == lexer.TokenOr
}

func (p *Parser) error(msg string) error {
	return fmt.Errorf("parser error at %d:%d: %s", p.current.Line, p.current.Column, msg)
}
