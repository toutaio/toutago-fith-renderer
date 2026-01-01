package runtime

import (
	"bytes"
	"fmt"

	"github.com/toutaio/toutago-fith-renderer/lexer"
	"github.com/toutaio/toutago-fith-renderer/parser"
)

// Runtime executes parsed templates with a given context.
type Runtime struct {
	context   *Context
	output    *bytes.Buffer
	functions *FunctionRegistry
}

// NewRuntime creates a new runtime with the given context.
func NewRuntime(ctx *Context) *Runtime {
	return &Runtime{
		context:   ctx,
		output:    &bytes.Buffer{},
		functions: NewFunctionRegistry(),
	}
}

// RegisterFunction adds a custom function to the runtime.
func (r *Runtime) RegisterFunction(name string, fn Function) {
	r.functions.Register(name, fn)
}

// GetContext returns the runtime's context.
func (r *Runtime) GetContext() *Context {
	return r.context
}

// ExecuteTemplate executes a template and stores the output.
func (r *Runtime) ExecuteTemplate(template *parser.Template) error {
	return r.executeTemplate(template)
}

// Output returns the rendered output.
func (r *Runtime) Output() string {
	return r.output.String()
}

// Execute executes a parsed template and returns the rendered output.
func Execute(template *parser.Template, ctx *Context) (string, error) {
	rt := NewRuntime(ctx)
	err := rt.executeTemplate(template)
	if err != nil {
		return "", err
	}
	return rt.output.String(), nil
}

// executeTemplate executes the template root node.
func (r *Runtime) executeTemplate(template *parser.Template) error {
	for _, node := range template.Nodes {
		if err := r.executeNode(node); err != nil {
			return err
		}
	}
	return nil
}

// executeNode executes a single AST node.
func (r *Runtime) executeNode(node parser.Node) error {
	switch n := node.(type) {
	case *parser.TextNode:
		return r.executeText(n)
	case *parser.VariableNode:
		return r.executeVariable(n)
	case *parser.IfNode:
		return r.executeIf(n)
	case *parser.RangeNode:
		return r.executeRange(n)
	case *parser.BinaryOpNode:
		val, err := r.evaluateBinaryOp(n)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprint(r.output, val)
		return nil
	case *parser.UnaryOpNode:
		val, err := r.evaluateUnaryOp(n)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprint(r.output, val)
		return nil
	case *parser.LiteralNode:
		val, err := r.evaluateLiteral(n)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprint(r.output, val)
		return nil
	case *parser.CallNode:
		return r.executeCall(n)
	case *parser.PipeNode:
		return r.executePipe(n)
	case *parser.IndexNode:
		val, err := r.evaluateIndex(n)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprint(r.output, val)
		return nil
	default:
		return fmt.Errorf("unsupported node type: %T", node)
	}
}

// executeText executes a text node by writing it to output.
func (r *Runtime) executeText(node *parser.TextNode) error {
	r.output.WriteString(node.Value)
	return nil
}

// executeVariable executes a variable node.
func (r *Runtime) executeVariable(node *parser.VariableNode) error {
	val, err := r.context.Get(node.Path)
	if err != nil {
		return fmt.Errorf("variable error at %d:%d: %v", node.Position.Line, node.Position.Column, err)
	}

	// Convert value to string and write to output
	_, _ = fmt.Fprint(r.output, val)
	return nil
}

// executeIf executes an if/else statement.
func (r *Runtime) executeIf(node *parser.IfNode) error {
	// Evaluate condition
	condVal, err := r.evaluateExpression(node.Condition)
	if err != nil {
		return fmt.Errorf("if condition error at %d:%d: %v", node.Position.Line, node.Position.Column, err)
	}

	// Check if condition is truthy
	if IsTruthy(condVal) {
		// Execute then branch
		for _, n := range node.Then {
			if err := r.executeNode(n); err != nil {
				return err
			}
		}
	} else if node.Else != nil {
		// Execute else branch
		for _, n := range node.Else {
			if err := r.executeNode(n); err != nil {
				return err
			}
		}
	}

	return nil
}

// executeRange executes a range loop.
func (r *Runtime) executeRange(node *parser.RangeNode) error {
	// Evaluate collection
	collVal, err := r.evaluateExpression(node.Collection)
	if err != nil {
		return fmt.Errorf("range collection error at %d:%d: %v", node.Position.Line, node.Position.Column, err)
	}

	// Try to convert to slice
	if items, ok := ToSlice(collVal); ok {
		return r.executeRangeSlice(node, items)
	}

	// Try to convert to map
	if keys, vals, ok := ToMap(collVal); ok {
		return r.executeRangeMap(node, keys, vals)
	}

	return fmt.Errorf("range error at %d:%d: value is not iterable", node.Position.Line, node.Position.Column)
}

// executeRangeSlice executes a range loop over a slice.
func (r *Runtime) executeRangeSlice(node *parser.RangeNode, items []interface{}) error {
	for idx, item := range items {
		// Push new scope
		r.context.PushScope()

		// Set loop variables
		r.context.Set(".", item) // Current item
		r.context.Set("@index", idx)
		r.context.Set("@first", idx == 0)
		r.context.Set("@last", idx == len(items)-1)

		// Execute loop body
		for _, n := range node.Body {
			if err := r.executeNode(n); err != nil {
				r.context.PopScope()
				return err
			}
		}

		// Pop scope
		r.context.PopScope()
	}

	return nil
}

// executeRangeMap executes a range loop over a map.
func (r *Runtime) executeRangeMap(
	node *parser.RangeNode,
	keys, vals []interface{},
) error {
	for idx, key := range keys {
		// Push new scope
		r.context.PushScope()

		// Set loop variables
		r.context.Set(".", vals[idx]) // Current value
		r.context.Set("@key", key)
		r.context.Set("@index", idx)
		r.context.Set("@first", idx == 0)
		r.context.Set("@last", idx == len(keys)-1)

		// Execute loop body
		for _, n := range node.Body {
			if err := r.executeNode(n); err != nil {
				r.context.PopScope()
				return err
			}
		}

		// Pop scope
		r.context.PopScope()
	}

	return nil
}

// executeCall executes a function call.
func (r *Runtime) executeCall(node *parser.CallNode) error {
	// Special case: if it's a no-arg call starting with @, treat it as a variable
	if len(node.Args) == 0 && node.Function != "" && node.Function[0] == '@' {
		// This is a special loop variable like @index, @first, @last
		val, err := r.context.Get([]string{node.Function})
		if err != nil {
			return fmt.Errorf("variable error at %d:%d: %v", node.Position.Line, node.Position.Column, err)
		}
		_, _ = fmt.Fprint(r.output, val)
		return nil
	}

	// Evaluate all arguments
	args := make([]interface{}, len(node.Args))
	for i, argNode := range node.Args {
		val, err := r.evaluateExpression(argNode)
		if err != nil {
			return fmt.Errorf("function argument error at %d:%d: %v", node.Position.Line, node.Position.Column, err)
		}
		args[i] = val
	}

	// Call the function
	result, err := r.functions.Call(node.Function, args...)
	if err != nil {
		return fmt.Errorf("function call error at %d:%d: %v", node.Position.Line, node.Position.Column, err)
	}

	// Output the result
	_, _ = fmt.Fprint(r.output, result)
	return nil
}

// executePipe executes a pipe expression.
func (r *Runtime) executePipe(node *parser.PipeNode) error {
	// Evaluate the initial value
	val, err := r.evaluateExpression(node.Value)
	if err != nil {
		return fmt.Errorf("pipe value error: %v", err)
	}

	// Apply each filter in sequence
	for _, filterName := range node.Filters {
		val, err = r.functions.Call(filterName, val)
		if err != nil {
			return fmt.Errorf("filter error (%s): %v", filterName, err)
		}
	}

	// Output the final result
	_, _ = fmt.Fprint(r.output, val)
	return nil
}

// evaluateExpression evaluates an expression node and returns its value.
func (r *Runtime) evaluateExpression(node parser.Node) (interface{}, error) {
	switch n := node.(type) {
	case *parser.VariableNode:
		return r.context.Get(n.Path)
	case *parser.LiteralNode:
		return r.evaluateLiteral(n)
	case *parser.BinaryOpNode:
		return r.evaluateBinaryOp(n)
	case *parser.UnaryOpNode:
		return r.evaluateUnaryOp(n)
	case *parser.IndexNode:
		return r.evaluateIndex(n)
	case *parser.CallNode:
		// Special case: @variables
		if len(n.Args) == 0 && n.Function != "" && n.Function[0] == '@' {
			return r.context.Get([]string{n.Function})
		}
		// Evaluate function calls
		args := make([]interface{}, len(n.Args))
		for i, argNode := range n.Args {
			val, err := r.evaluateExpression(argNode)
			if err != nil {
				return nil, err
			}
			args[i] = val
		}
		return r.functions.Call(n.Function, args...)
	default:
		return nil, fmt.Errorf("cannot evaluate node type: %T", node)
	}
}

// evaluateLiteral evaluates a literal node.
func (r *Runtime) evaluateLiteral(node *parser.LiteralNode) (interface{}, error) {
	return node.Value, nil
}

// evaluateBinaryOp evaluates a binary operation.
func (r *Runtime) evaluateBinaryOp(node *parser.BinaryOpNode) (interface{}, error) {
	left, err := r.evaluateExpression(node.Left)
	if err != nil {
		return nil, err
	}

	right, err := r.evaluateExpression(node.Right)
	if err != nil {
		return nil, err
	}

	return r.applyBinaryOperator(node.Operator, left, right)
}

// applyBinaryOperator applies a binary operator to two values.
func (r *Runtime) applyBinaryOperator(op lexer.TokenType, left, right interface{}) (interface{}, error) {
	// Comparison operators
	if result, ok := r.tryComparisonOp(op, left, right); ok {
		return result, nil
	}

	// Logical operators
	if result, ok := r.tryLogicalOp(op, left, right); ok {
		return result, nil
	}

	// Arithmetic operators
	return r.tryArithmeticOp(op, left, right)
}

// tryComparisonOp attempts to apply a comparison operator.
func (r *Runtime) tryComparisonOp(op lexer.TokenType, left, right interface{}) (interface{}, bool) {
	switch op {
	case lexer.TokenEqual:
		return compareEqual(left, right), true
	case lexer.TokenNotEqual:
		return !compareEqual(left, right), true
	case lexer.TokenLess:
		return compareLess(left, right), true
	case lexer.TokenGreater:
		return compareLess(right, left), true
	case lexer.TokenLessEq:
		return !compareLess(right, left), true
	case lexer.TokenGreaterEq:
		return !compareLess(left, right), true
	}
	return nil, false
}

// tryLogicalOp attempts to apply a logical operator.
func (r *Runtime) tryLogicalOp(op lexer.TokenType, left, right interface{}) (interface{}, bool) {
	switch op {
	case lexer.TokenAnd:
		return IsTruthy(left) && IsTruthy(right), true
	case lexer.TokenOr:
		return IsTruthy(left) || IsTruthy(right), true
	}
	return nil, false
}

// tryArithmeticOp attempts to apply an arithmetic operator.
func (r *Runtime) tryArithmeticOp(op lexer.TokenType, left, right interface{}) (interface{}, error) {
	switch op {
	case lexer.TokenPlus:
		return addValues(left, right)
	case lexer.TokenMinus:
		return subtractValues(left, right)
	case lexer.TokenMult:
		return multiplyValues(left, right)
	case lexer.TokenDiv:
		return divideValues(left, right)
	case lexer.TokenMod:
		return modValues(left, right)
	default:
		return nil, fmt.Errorf("unsupported operator: %v", op)
	}
}

// evaluateUnaryOp evaluates a unary operation.
func (r *Runtime) evaluateUnaryOp(node *parser.UnaryOpNode) (interface{}, error) {
	operand, err := r.evaluateExpression(node.Operand)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case lexer.TokenNot:
		return !IsTruthy(operand), nil
	case lexer.TokenMinus:
		return negateValue(operand)
	default:
		return nil, fmt.Errorf("unsupported unary operator: %v", node.Operator)
	}
}

// evaluateIndex evaluates an index expression.
func (r *Runtime) evaluateIndex(node *parser.IndexNode) (interface{}, error) {
	obj, err := r.evaluateExpression(node.Object)
	if err != nil {
		return nil, err
	}

	idx, err := r.evaluateExpression(node.Index)
	if err != nil {
		return nil, err
	}

	return r.context.GetIndex(obj, idx)
}

// Helper functions for arithmetic and comparison

func compareEqual(a, b interface{}) bool {
	return fmt.Sprint(a) == fmt.Sprint(b)
}

func compareLess(a, b interface{}) bool {
	aNum, aOk := toFloat(a)
	bNum, bOk := toFloat(b)
	if aOk && bOk {
		return aNum < bNum
	}
	return fmt.Sprint(a) < fmt.Sprint(b)
}

func toFloat(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	case float32:
		return float64(v), true
	default:
		return 0, false
	}
}

func toInt(val interface{}) (int, bool) {
	switch v := val.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	default:
		return 0, false
	}
}

func addValues(a, b interface{}) (interface{}, error) {
	aNum, aOk := toFloat(a)
	bNum, bOk := toFloat(b)
	if aOk && bOk {
		// Check if both were integers
		if aInt, aIsInt := toInt(a); aIsInt {
			if bInt, bIsInt := toInt(b); bIsInt {
				return aInt + bInt, nil
			}
		}
		return aNum + bNum, nil
	}
	return nil, fmt.Errorf("cannot add %T and %T", a, b)
}

func subtractValues(a, b interface{}) (interface{}, error) {
	aNum, aOk := toFloat(a)
	bNum, bOk := toFloat(b)
	if aOk && bOk {
		if aInt, aIsInt := toInt(a); aIsInt {
			if bInt, bIsInt := toInt(b); bIsInt {
				return aInt - bInt, nil
			}
		}
		return aNum - bNum, nil
	}
	return nil, fmt.Errorf("cannot subtract %T and %T", a, b)
}

func multiplyValues(a, b interface{}) (interface{}, error) {
	aNum, aOk := toFloat(a)
	bNum, bOk := toFloat(b)
	if aOk && bOk {
		if aInt, aIsInt := toInt(a); aIsInt {
			if bInt, bIsInt := toInt(b); bIsInt {
				return aInt * bInt, nil
			}
		}
		return aNum * bNum, nil
	}
	return nil, fmt.Errorf("cannot multiply %T and %T", a, b)
}

func divideValues(a, b interface{}) (interface{}, error) {
	aNum, aOk := toFloat(a)
	bNum, bOk := toFloat(b)
	if aOk && bOk {
		if bNum == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return aNum / bNum, nil
	}
	return nil, fmt.Errorf("cannot divide %T and %T", a, b)
}

func modValues(a, b interface{}) (interface{}, error) {
	aInt, aOk := toInt(a)
	bInt, bOk := toInt(b)
	if aOk && bOk {
		if bInt == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}
		return aInt % bInt, nil
	}
	return nil, fmt.Errorf("cannot mod %T and %T", a, b)
}

func negateValue(val interface{}) (interface{}, error) {
	if num, ok := toFloat(val); ok {
		if intVal, isInt := toInt(val); isInt {
			return -intVal, nil
		}
		return -num, nil
	}
	return nil, fmt.Errorf("cannot negate %T", val)
}
