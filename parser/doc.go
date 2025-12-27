// Package parser provides AST generation for Fíth templates.
//
// The parser consumes tokens from the lexer and builds an Abstract
// Syntax Tree (AST) that represents the structure of the template.
// The AST can then be compiled and executed by the runtime.
//
// Example:
//
//	import (
//	    "github.com/toutaio/toutago-fith-renderer/lexer"
//	    "github.com/toutaio/toutago-fith-renderer/parser"
//	)
//
//	l := lexer.New("Hello {{.Name}}")
//	p := parser.New(l)
//	ast, err := p.Parse()
//	if err != nil {
//	    log.Fatal(err)
//	}
package parser
