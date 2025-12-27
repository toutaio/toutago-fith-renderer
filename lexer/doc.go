// Package lexer provides tokenization for Fíth templates.
//
// The lexer scans template input and produces a stream of tokens
// that can be consumed by the parser. It handles both text content
// and template expressions delimited by {{ and }}.
//
// Example:
//
//	l := lexer.New("Hello {{.Name}}")
//	for {
//	    tok, err := l.NextToken()
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    if tok.Type == lexer.TokenEOF {
//	        break
//	    }
//	    fmt.Printf("%v\n", tok)
//	}
package lexer
