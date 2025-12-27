package lexer

import "fmt"

// TokenType represents the type of a lexical token.
type TokenType int

// Token types recognized by the lexer.
const (
	TokenError TokenType = iota // Error token
	TokenEOF                    // End of file

	// Literals
	TokenText   // Text content outside template expressions
	TokenIdent  // Identifier (variable name, function name)
	TokenString // String literal "..."
	TokenNumber // Number literal (int or float)

	// Delimiters
	TokenOpenDelim  // {{
	TokenCloseDelim // }}

	// Operators
	TokenDot       // .
	TokenPipe      // |
	TokenEqual     // ==
	TokenNotEqual  // !=
	TokenLess      // <
	TokenGreater   // >
	TokenLessEq    // <=
	TokenGreaterEq // >=
	TokenAnd       // &&
	TokenOr        // ||
	TokenNot       // !
	TokenPlus      // +
	TokenMinus     // -
	TokenMult      // *
	TokenDiv       // /
	TokenMod       // %
	TokenAssign    // =
	TokenComma     // ,
	TokenColon     // :
	TokenLParen    // (
	TokenRParen    // )
	TokenLBrack    // [
	TokenRBrack    // ]

	// Keywords
	TokenIf      // if
	TokenElse    // else
	TokenEnd     // end
	TokenRange   // range
	TokenInclude // include
	TokenExtends // extends
	TokenBlock   // block
)

// String returns the string representation of the token type.
func (t TokenType) String() string {
	names := map[TokenType]string{
		TokenError:      "ERROR",
		TokenEOF:        "EOF",
		TokenText:       "TEXT",
		TokenIdent:      "IDENT",
		TokenString:     "STRING",
		TokenNumber:     "NUMBER",
		TokenOpenDelim:  "{{",
		TokenCloseDelim: "}}",
		TokenDot:        ".",
		TokenPipe:       "|",
		TokenEqual:      "==",
		TokenNotEqual:   "!=",
		TokenLess:       "<",
		TokenGreater:    ">",
		TokenLessEq:     "<=",
		TokenGreaterEq:  ">=",
		TokenAnd:        "&&",
		TokenOr:         "||",
		TokenNot:        "!",
		TokenPlus:       "+",
		TokenMinus:      "-",
		TokenMult:       "*",
		TokenDiv:        "/",
		TokenMod:        "%",
		TokenAssign:     "=",
		TokenComma:      ",",
		TokenColon:      ":",
		TokenLParen:     "(",
		TokenRParen:     ")",
		TokenLBrack:     "[",
		TokenRBrack:     "]",
		TokenIf:         "IF",
		TokenElse:       "ELSE",
		TokenEnd:        "END",
		TokenRange:      "RANGE",
		TokenInclude:    "INCLUDE",
		TokenExtends:    "EXTENDS",
		TokenBlock:      "BLOCK",
	}
	if name, ok := names[t]; ok {
		return name
	}
	return fmt.Sprintf("TokenType(%d)", t)
}

// Token represents a lexical token.
type Token struct {
	Type   TokenType // Type of token
	Value  string    // Literal value of token
	Line   int       // Line number (1-indexed)
	Column int       // Column number (1-indexed)
}

// String returns a string representation of the token for debugging.
func (t Token) String() string {
	if len(t.Value) > 20 {
		return fmt.Sprintf("%s(%q...) at %d:%d", t.Type, t.Value[:20], t.Line, t.Column)
	}
	return fmt.Sprintf("%s(%q) at %d:%d", t.Type, t.Value, t.Line, t.Column)
}

// keywords maps keyword strings to their token types.
var keywords = map[string]TokenType{
	"if":      TokenIf,
	"else":    TokenElse,
	"end":     TokenEnd,
	"range":   TokenRange,
	"include": TokenInclude,
	"extends": TokenExtends,
	"block":   TokenBlock,
}

// IsKeyword checks if a string is a keyword and returns its TokenType.
func IsKeyword(ident string) (TokenType, bool) {
	tok, ok := keywords[ident]
	return tok, ok
}
