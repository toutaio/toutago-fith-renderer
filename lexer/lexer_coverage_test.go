package lexer

import (
	"testing"
)

// Additional tests for improved coverage

func TestLexer_ErrorToken(t *testing.T) {
	l := New("{{&}}")
	tokens, err := l.All()
	if err == nil {
		t.Error("expected error for single &, got nil")
	}
	if len(tokens) == 0 {
		t.Error("expected at least opening delimiter token")
	}
}

func TestLexer_AllMethod(t *testing.T) {
	l := New("{{.x}}")
	tokens, err := l.All()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) < 4 {
		t.Errorf("expected at least 4 tokens, got %d", len(tokens))
	}
}

func TestLexer_PeekAhead(t *testing.T) {
	l := New("abc")
	if l.peekAhead(10) != 0 {
		t.Error("peekAhead beyond bounds should return 0")
	}
}

func TestLexer_MultipleExpressions_Extended(t *testing.T) {
	l := New("{{.a}}{{.b}}")
	tokens, err := l.All()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Count identifiers
	count := 0
	for _, tok := range tokens {
		if tok.Type == TokenIdent {
			count++
		}
	}
	if count != 2 {
		t.Errorf("expected 2 identifiers, got %d", count)
	}
}

func TestLexer_StringWithEscapes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`{{"hello\"world"}}`, `hello\"world`},
		{`{{'it\'s'}}`, `it\'s`},
		{`{{"test\\path"}}`, `test\\path`},
	}

	for _, tt := range tests {
		l := New(tt.input)
		_, _ = l.NextToken() // {{
		tok, err := l.NextToken()
		if err != nil {
			t.Errorf("input %q: unexpected error: %v", tt.input, err)
			continue
		}
		if tok.Type != TokenString {
			t.Errorf("input %q: expected TokenString, got %v", tt.input, tok.Type)
		}
		if tok.Value != tt.want {
			t.Errorf("input %q: expected %q, got %q", tt.input, tt.want, tok.Value)
		}
	}
}

func TestLexer_NumberFormats(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"{{42}}", "42"},
		{"{{3.14}}", "3.14"},
		{"{{0.5}}", "0.5"},
		{"{{100.0}}", "100.0"},
	}

	for _, tt := range tests {
		l := New(tt.input)
		_, _ = l.NextToken() // {{
		tok, err := l.NextToken()
		if err != nil {
			t.Errorf("input %q: unexpected error: %v", tt.input, err)
			continue
		}
		if tok.Type != TokenNumber {
			t.Errorf("input %q: expected TokenNumber, got %v", tt.input, tok.Type)
		}
		if tok.Value != tt.want {
			t.Errorf("input %q: expected %q, got %q", tt.input, tt.want, tok.Value)
		}
	}
}

func TestLexer_IdentifierStartingWithAt(t *testing.T) {
	l := New("{{@index}}")
	_, _ = l.NextToken() // {{
	tok, err := l.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != TokenIdent {
		t.Errorf("expected TokenIdent, got %v", tok.Type)
	}
	if tok.Value != "@index" {
		t.Errorf("expected '@index', got %q", tok.Value)
	}
}

func TestLexer_UnterminatedString(t *testing.T) {
	inputs := []string{
		`{{"unclosed`,
		`{{'unclosed`,
	}

	for _, input := range inputs {
		l := New(input)
		_, _ = l.NextToken() // {{
		_, err := l.NextToken()
		if err == nil {
			t.Errorf("input %q: expected error for unterminated/invalid string", input)
		}
	}
}

func TestLexer_AllOperators(t *testing.T) {
	operators := []struct {
		input string
		typ   TokenType
	}{
		{"{{+}}", TokenPlus},
		{"{{-}}", TokenMinus},
		{"{{*}}", TokenMult},
		{"{{/}}", TokenDiv},
		{"{{%}}", TokenMod},
		{"{{==}}", TokenEqual},
		{"{{!=}}", TokenNotEqual},
		{"{{<}}", TokenLess},
		{"{{>}}", TokenGreater},
		{"{{<=}}", TokenLessEq},
		{"{{>=}}", TokenGreaterEq},
		{"{{&&}}", TokenAnd},
		{"{{||}}", TokenOr},
		{"{{|}}", TokenPipe},
		{"{{!}}", TokenNot},
		{"{{=}}", TokenAssign},
	}

	for _, tt := range operators {
		l := New(tt.input)
		_, _ = l.NextToken() // {{
		tok, err := l.NextToken()
		if err != nil {
			t.Errorf("input %q: unexpected error: %v", tt.input, err)
			continue
		}
		if tok.Type != tt.typ {
			t.Errorf("input %q: expected %v, got %v", tt.input, tt.typ, tok.Type)
		}
	}
}

func TestLexer_Keywords_Extended(t *testing.T) {
	keywords := map[string]TokenType{
		"if":      TokenIf,
		"else":    TokenElse,
		"end":     TokenEnd,
		"range":   TokenRange,
		"include": TokenInclude,
		"extends": TokenExtends,
		"block":   TokenBlock,
	}

	for word, expectedType := range keywords {
		input := "{{" + word + "}}"
		l := New(input)
		_, _ = l.NextToken() // {{
		tok, err := l.NextToken()
		if err != nil {
			t.Errorf("input %q: unexpected error: %v", word, err)
			continue
		}
		if tok.Type != expectedType {
			t.Errorf("word %q: expected %v, got %v", word, expectedType, tok.Type)
		}
	}
}
