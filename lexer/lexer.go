package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// Lexer tokenizes template input into a stream of tokens.
// It maintains position information for error reporting.
type Lexer struct {
	input     string // The input string being lexed
	pos       int    // Current position in input (bytes)
	line      int    // Current line number (1-indexed)
	column    int    // Current column number (1-indexed)
	start     int    // Start position of current token
	startLine int    // Line number at token start
	startCol  int    // Column number at token start
	inExpr    bool   // True if inside template expression {{...}}
}

// New creates a new Lexer for the given input string.
func New(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
	}
}

// NextToken returns the next token from the input.
// Returns TokenEOF when the end of input is reached.
// Returns a token with TokenError type if invalid syntax is encountered.
func (l *Lexer) NextToken() (Token, error) {
	// Skip whitespace if inside expression
	if l.inExpr {
		l.skipWhitespace()
	}

	// Mark start of new token
	l.start = l.pos
	l.startLine = l.line
	l.startCol = l.column

	// Check for EOF
	if l.pos >= len(l.input) {
		return l.makeToken(TokenEOF, ""), nil
	}

	// If not in expression, scan text until we hit {{
	if !l.inExpr {
		return l.scanText()
	}

	// Inside expression - scan tokens
	return l.scanExpression()
}

// scanText scans literal text until we encounter {{.
func (l *Lexer) scanText() (Token, error) {
	start := l.pos

	for l.pos < len(l.input) {
		// Look for opening delimiter
		if l.pos+1 < len(l.input) && l.input[l.pos] == '{' && l.input[l.pos+1] == '{' {
			// Found opening delimiter
			if l.pos > start {
				// Return text token before delimiter
				text := l.input[start:l.pos]
				return l.makeToken(TokenText, text), nil
			}
			// No text before delimiter, switch to expression mode
			l.inExpr = true
			l.advance()
			l.advance()
			return l.makeToken(TokenOpenDelim, "{{"), nil
		}

		l.advance()
	}

	// Reached EOF while scanning text
	if l.pos > start {
		text := l.input[start:l.pos]
		return l.makeToken(TokenText, text), nil
	}

	return l.makeToken(TokenEOF, ""), nil
}

// scanExpression scans a single token inside a template expression.
func (l *Lexer) scanExpression() (Token, error) {
	ch := l.peek()

	// Check for closing delimiter
	if ch == '}' && l.peekAhead(1) == '}' {
		l.advance()
		l.advance()
		l.inExpr = false
		return l.makeToken(TokenCloseDelim, "}}"), nil
	}

	// Try single character tokens
	if tok, ok := l.trySingleCharToken(ch); ok {
		return tok, nil
	}

	// Try multi-character operators
	if tok, err := l.tryMultiCharOperator(ch); tok.Type != "" || err != nil {
		return tok, err
	}

	// String literals
	if ch == '"' || ch == '\'' {
		return l.scanString(ch)
	}

	// Numbers
	if unicode.IsDigit(rune(ch)) {
		return l.scanNumber()
	}

	// Identifiers and keywords
	if unicode.IsLetter(rune(ch)) || ch == '_' || ch == '@' {
		return l.scanIdentifier()
	}

	// Unknown character
	return l.errorToken(fmt.Sprintf("unexpected character: %q", ch))
}

// trySingleCharToken attempts to scan a single character token.
func (l *Lexer) trySingleCharToken(ch byte) (Token, bool) {
	var tokType TokenType
	var lexeme string

	switch ch {
	case '.':
		tokType, lexeme = TokenDot, "."
	case '+':
		tokType, lexeme = TokenPlus, "+"
	case '-':
		tokType, lexeme = TokenMinus, "-"
	case '*':
		tokType, lexeme = TokenMult, "*"
	case '/':
		tokType, lexeme = TokenDiv, "/"
	case '%':
		tokType, lexeme = TokenMod, "%"
	case ',':
		tokType, lexeme = TokenComma, ","
	case ':':
		tokType, lexeme = TokenColon, ":"
	case '(':
		tokType, lexeme = TokenLParen, "("
	case ')':
		tokType, lexeme = TokenRParen, ")"
	case '[':
		tokType, lexeme = TokenLBrack, "["
	case ']':
		tokType, lexeme = TokenRBrack, "]"
	default:
		return Token{}, false
	}

	l.advance()
	return l.makeToken(tokType, lexeme), true
}

// tryMultiCharOperator attempts to scan multi-character operators.
func (l *Lexer) tryMultiCharOperator(ch byte) (Token, error) {
	switch ch {
	case '=':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return l.makeToken(TokenEqual, "=="), nil
		}
		return l.makeToken(TokenAssign, "="), nil

	case '!':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return l.makeToken(TokenNotEqual, "!="), nil
		}
		return l.makeToken(TokenNot, "!"), nil

	case '<':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return l.makeToken(TokenLessEq, "<="), nil
		}
		return l.makeToken(TokenLess, "<"), nil

	case '>':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return l.makeToken(TokenGreaterEq, ">="), nil
		}
		return l.makeToken(TokenGreater, ">"), nil

	case '&':
		l.advance()
		if l.peek() == '&' {
			l.advance()
			return l.makeToken(TokenAnd, "&&"), nil
		}
		return Token{}, l.errorToken("expected '&&', got '&'")

	case '|':
		l.advance()
		if l.peek() == '|' {
			l.advance()
			return l.makeToken(TokenOr, "||"), nil
		}
		return l.makeToken(TokenPipe, "|"), nil
	}

	return Token{}, nil
}

// scanString scans a string literal.
func (l *Lexer) scanString(quote byte) (Token, error) {
	start := l.pos
	l.advance() // Skip opening quote

	for l.pos < len(l.input) {
		ch := l.input[l.pos]

		if ch == quote {
			l.advance()                         // Skip closing quote
			value := l.input[start+1 : l.pos-1] // Exclude quotes
			return l.makeToken(TokenString, value), nil
		}

		if ch == '\\' {
			l.advance() // Skip escape character
			if l.pos < len(l.input) {
				l.advance() // Skip escaped character
			}
			continue
		}

		if ch == '\n' {
			return l.errorToken("unterminated string literal")
		}

		l.advance()
	}

	return l.errorToken("unterminated string literal")
}

// scanNumber scans a number literal (integer or float).
func (l *Lexer) scanNumber() (Token, error) {
	start := l.pos

	// Scan digits
	for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
		l.advance()
	}

	// Check for decimal point
	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		// Look ahead to see if next char is a digit
		if l.pos+1 < len(l.input) && unicode.IsDigit(rune(l.input[l.pos+1])) {
			l.advance() // Skip .
			for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
				l.advance()
			}
		}
	}

	value := l.input[start:l.pos]
	return l.makeToken(TokenNumber, value), nil
}

// scanIdentifier scans an identifier or keyword.
func (l *Lexer) scanIdentifier() (Token, error) {
	start := l.pos

	// First character already validated (letter, _, or @)
	l.advance()

	// Continue with letters, digits, or underscores
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' {
			l.advance()
		} else {
			break
		}
	}

	value := l.input[start:l.pos]

	// Check if it's a keyword
	if tokType, ok := IsKeyword(value); ok {
		return l.makeToken(tokType, value), nil
	}

	return l.makeToken(TokenIdent, value), nil
}

// skipWhitespace skips whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			if ch == '\n' {
				l.line++
				l.column = 0
			}
			l.advance()
		} else {
			break
		}
	}
}

// peek returns the current character without advancing.
func (l *Lexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

// peekAhead returns the character n positions ahead without advancing.
func (l *Lexer) peekAhead(n int) byte {
	pos := l.pos + n
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
}

// advance moves to the next character.
func (l *Lexer) advance() {
	if l.pos < len(l.input) {
		r, size := utf8.DecodeRuneInString(l.input[l.pos:])
		l.pos += size
		if r == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
	}
}

// makeToken creates a token with current position information.
func (l *Lexer) makeToken(typ TokenType, value string) Token {
	return Token{
		Type:   typ,
		Value:  value,
		Line:   l.startLine,
		Column: l.startCol,
	}
}

// errorToken creates an error token with current position.
func (l *Lexer) errorToken(message string) (Token, error) {
	tok := l.makeToken(TokenError, message)
	return tok, fmt.Errorf("lexer error at %d:%d: %s", tok.Line, tok.Column, message)
}

// All returns all tokens from the input.
// This is a convenience method for testing and debugging.
func (l *Lexer) All() ([]Token, error) {
	var tokens []Token
	for {
		tok, err := l.NextToken()
		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}
	return tokens, nil
}
