// Package runtime provides the execution engine for Fíth templates.
//
// The runtime package takes parsed AST from the parser and executes it
// with a given data context, producing the final rendered output.
//
// Example:
//
//	import (
//	    "github.com/toutaio/toutago-fith-renderer/lexer"
//	    "github.com/toutaio/toutago-fith-renderer/parser"
//	    "github.com/toutaio/toutago-fith-renderer/runtime"
//	)
//
//	l := lexer.New("Hello {{.Name}}")
//	p := parser.New(l)
//	ast, _ := p.Parse()
//
//	ctx := runtime.NewContext(map[string]interface{}{"Name": "World"})
//	output, _ := runtime.Execute(ast, ctx)
//	fmt.Println(output) // "Hello World"
package runtime
