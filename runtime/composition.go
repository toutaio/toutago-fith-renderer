package runtime

import (
	"fmt"

	"github.com/toutaio/toutago-fith-renderer/parser"
)

// Loader is the interface for loading templates.
// The runtime uses this to load included and extended templates.
type Loader interface {
	Load(slug string) (*parser.Template, error)
	Exists(slug string) bool
}

// CompositionRuntime extends Runtime with template composition capabilities.
type CompositionRuntime struct {
	*Runtime
	loader          Loader
	blocks          map[string][]parser.Node
	includeStack    []string
	maxIncludeDepth int
}

// NewCompositionRuntime creates a runtime with composition support.
func NewCompositionRuntime(ctx *Context, loader Loader) *CompositionRuntime {
	return &CompositionRuntime{
		Runtime:         NewRuntime(ctx),
		loader:          loader,
		blocks:          make(map[string][]parser.Node),
		includeStack:    make([]string, 0),
		maxIncludeDepth: 100, // Prevent deep recursion
	}
}

// ExecuteWithLoader executes a template with loader support for composition.
func ExecuteWithLoader(template *parser.Template, ctx *Context, loader Loader) (string, error) {
	rt := NewCompositionRuntime(ctx, loader)

	// Check if template has extends directive
	extendsNode := rt.findExtendsNode(template)
	if extendsNode != nil {
		return rt.executeWithExtends(template, extendsNode)
	}

	// Execute normally with composition support
	for _, node := range template.Nodes {
		if err := rt.executeNode(node); err != nil {
			return "", err
		}
	}
	return rt.output.String(), nil
}

// findExtendsNode finds the extends directive in a template (must be first non-text node).
func (r *CompositionRuntime) findExtendsNode(tmpl *parser.Template) *parser.ExtendsNode {
	for _, node := range tmpl.Nodes {
		// Skip leading whitespace text nodes
		if textNode, ok := node.(*parser.TextNode); ok {
			if isWhitespace(textNode.Value) {
				continue
			}
			// Non-whitespace text before extends is an error
			return nil
		}

		if extendsNode, ok := node.(*parser.ExtendsNode); ok {
			return extendsNode
		}

		// extends must come first
		return nil
	}
	return nil
}

// executeWithExtends handles template inheritance.
func (r *CompositionRuntime) executeWithExtends(child *parser.Template, extendsNode *parser.ExtendsNode) (string, error) {
	// Collect blocks from child template
	r.collectBlocks(child)

	// Load parent template
	parent, err := r.loader.Load(extendsNode.Template)
	if err != nil {
		return "", fmt.Errorf("failed to load parent template %q: %w", extendsNode.Template, err)
	}

	// Check if parent also extends
	parentExtendsNode := r.findExtendsNode(parent)
	if parentExtendsNode != nil {
		// Recursive extends
		return r.executeWithExtends(parent, parentExtendsNode)
	}

	// Execute parent template with child blocks
	for _, node := range parent.Nodes {
		if err := r.executeNode(node); err != nil {
			return "", err
		}
	}
	return r.output.String(), nil
}

// collectBlocks collects all block definitions from a template.
func (r *CompositionRuntime) collectBlocks(tmpl *parser.Template) {
	for _, node := range tmpl.Nodes {
		if blockNode, ok := node.(*parser.BlockNode); ok {
			r.blocks[blockNode.Name] = blockNode.Body
		}
	}
}

// executeNode overrides the base executeNode to handle composition nodes.
func (r *CompositionRuntime) executeNode(node parser.Node) error {
	switch n := node.(type) {
	case *parser.IncludeNode:
		return r.executeInclude(n)
	case *parser.BlockNode:
		return r.executeBlock(n)
	case *parser.ExtendsNode:
		// Extends handled separately, skip here
		return nil
	default:
		// Delegate to base runtime
		return r.Runtime.executeNode(node)
	}
}

// executeInclude handles template inclusion.
func (r *CompositionRuntime) executeInclude(node *parser.IncludeNode) error {
	// Check for circular includes
	for _, slug := range r.includeStack {
		if slug == node.Template {
			return fmt.Errorf("circular include detected: %q", node.Template)
		}
	}

	// Check depth limit
	if len(r.includeStack) >= r.maxIncludeDepth {
		return fmt.Errorf("maximum include depth exceeded (%d)", r.maxIncludeDepth)
	}

	// Load the included template
	tmpl, err := r.loader.Load(node.Template)
	if err != nil {
		return fmt.Errorf("failed to load include %q: %w", node.Template, err)
	}

	// Push to include stack
	r.includeStack = append(r.includeStack, node.Template)
	defer func() {
		// Pop from include stack
		r.includeStack = r.includeStack[:len(r.includeStack)-1]
	}()

	// Create new context for include
	includeCtx := r.createIncludeContext(node)

	// Save and restore context
	savedCtx := r.context
	r.context = includeCtx
	defer func() {
		r.context = savedCtx
	}()

	// Execute included template
	for _, node := range tmpl.Nodes {
		if err := r.executeNode(node); err != nil {
			return err
		}
	}
	return nil
}

// createIncludeContext creates a context for an included template.
func (r *CompositionRuntime) createIncludeContext(node *parser.IncludeNode) *Context {
	// If context is explicitly provided, use it
	if node.Context != nil {
		val, err := r.Runtime.evaluateExpression(node.Context)
		if err != nil {
			// Fall back to current context on error
			return r.context
		}
		return NewContext(val)
	}

	// If parameters are provided, create new context with params
	if len(node.Params) > 0 {
		data := make(map[string]interface{})
		for key, valueNode := range node.Params {
			val, err := r.Runtime.evaluateExpression(valueNode)
			if err != nil {
				// Skip invalid params
				continue
			}
			data[key] = val
		}
		return NewContext(data)
	}

	// Otherwise inherit current context
	return r.context
}

// executeBlock handles block directives.
func (r *CompositionRuntime) executeBlock(node *parser.BlockNode) error {
	// Check if block has been overridden
	if overrideBody, exists := r.blocks[node.Name]; exists {
		// Execute override content
		for _, n := range overrideBody {
			if err := r.executeNode(n); err != nil {
				return err
			}
		}
		return nil
	}

	// Execute default block content
	for _, n := range node.Body {
		if err := r.executeNode(n); err != nil {
			return err
		}
	}
	return nil
}

// isWhitespace checks if a string contains only whitespace.
func isWhitespace(s string) bool {
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			return false
		}
	}
	return true
}
