package lexer

import (
	"testing"
)

func TestLexer_SimpleText(t *testing.T) {
	input := "Hello, World!"
	l := New(input)

	tok, err := l.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tok.Type != TokenText {
		t.Errorf("expected TokenText, got %v", tok.Type)
	}

	if tok.Value != input {
		t.Errorf("expected %q, got %q", input, tok.Value)
	}

	// Should get EOF next
	tok, err = l.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tok.Type != TokenEOF {
		t.Errorf("expected TokenEOF, got %v", tok.Type)
	}
}

func TestLexer_SimpleVariable(t *testing.T) {
	input := "{{.Name}}"
	l := New(input)

	tests := []struct {
		wantType  TokenType
		wantValue string
	}{
		{TokenOpenDelim, "{{"},
		{TokenDot, "."},
		{TokenIdent, "Name"},
		{TokenCloseDelim, "}}"},
		{TokenEOF, ""},
	}

	for i, tt := range tests {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatalf("test %d: unexpected error: %v", i, err)
		}

		if tok.Type != tt.wantType {
			t.Errorf("test %d: expected type %v, got %v", i, tt.wantType, tok.Type)
		}

		if tok.Value != tt.wantValue {
			t.Errorf("test %d: expected value %q, got %q", i, tt.wantValue, tok.Value)
		}
	}
}

func TestLexer_MixedTextAndExpressions(t *testing.T) {
	input := "Hello {{.Name}}!"
	l := New(input)

	tests := []struct {
		wantType  TokenType
		wantValue string
	}{
		{TokenText, "Hello "},
		{TokenOpenDelim, "{{"},
		{TokenDot, "."},
		{TokenIdent, "Name"},
		{TokenCloseDelim, "}}"},
		{TokenText, "!"},
		{TokenEOF, ""},
	}

	for i, tt := range tests {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatalf("test %d: unexpected error: %v", i, err)
		}

		if tok.Type != tt.wantType {
			t.Errorf("test %d: expected type %v, got %v", i, tt.wantType, tok.Type)
		}

		if tok.Value != tt.wantValue {
			t.Errorf("test %d: expected value %q, got %q", i, tt.wantValue, tok.Value)
		}
	}
}

func TestLexer_NestedFieldAccess(t *testing.T) {
	input := "{{.User.Profile.Name}}"
	l := New(input)

	tests := []struct {
		wantType  TokenType
		wantValue string
	}{
		{TokenOpenDelim, "{{"},
		{TokenDot, "."},
		{TokenIdent, "User"},
		{TokenDot, "."},
		{TokenIdent, "Profile"},
		{TokenDot, "."},
		{TokenIdent, "Name"},
		{TokenCloseDelim, "}}"},
		{TokenEOF, ""},
	}

	for i, tt := range tests {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatalf("test %d: unexpected error: %v", i, err)
		}

		if tok.Type != tt.wantType {
			t.Errorf("test %d: expected type %v, got %v", i, tt.wantType, tok.Type)
		}
	}
}

func TestLexer_Keywords(t *testing.T) {
	tests := []struct {
		input string
		want  TokenType
	}{
		{"{{if .Active}}", TokenIf},
		{"{{else}}", TokenElse},
		{"{{end}}", TokenEnd},
		{"{{range .Items}}", TokenRange},
		{"{{include \"header\"}}", TokenInclude},
		{"{{extends \"layout\"}}", TokenExtends},
		{"{{block \"content\"}}", TokenBlock},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(tt.input)
			_, _ = l.NextToken() // Skip {{

			tok, err := l.NextToken()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tok.Type != tt.want {
				t.Errorf("expected %v, got %v", tt.want, tok.Type)
			}
		})
	}
}

func TestLexer_Operators(t *testing.T) {
	tests := []struct {
		input string
		want  []TokenType
	}{
		{
			"{{.A == .B}}",
			[]TokenType{TokenOpenDelim, TokenDot, TokenIdent, TokenEqual, TokenDot, TokenIdent, TokenCloseDelim},
		},
		{
			"{{.A != .B}}",
			[]TokenType{TokenOpenDelim, TokenDot, TokenIdent, TokenNotEqual, TokenDot, TokenIdent, TokenCloseDelim},
		},
		{"{{.A < .B}}", []TokenType{TokenOpenDelim, TokenDot, TokenIdent, TokenLess, TokenDot, TokenIdent, TokenCloseDelim}},
		{
			"{{.A > .B}}",
			[]TokenType{TokenOpenDelim, TokenDot, TokenIdent, TokenGreater, TokenDot, TokenIdent, TokenCloseDelim},
		},
		{
			"{{.A <= .B}}",
			[]TokenType{TokenOpenDelim, TokenDot, TokenIdent, TokenLessEq, TokenDot, TokenIdent, TokenCloseDelim},
		},
		{
			"{{.A >= .B}}",
			[]TokenType{TokenOpenDelim, TokenDot, TokenIdent, TokenGreaterEq, TokenDot, TokenIdent, TokenCloseDelim},
		},
		{"{{.A && .B}}", []TokenType{TokenOpenDelim, TokenDot, TokenIdent, TokenAnd, TokenDot, TokenIdent, TokenCloseDelim}},
		{"{{.A || .B}}", []TokenType{TokenOpenDelim, TokenDot, TokenIdent, TokenOr, TokenDot, TokenIdent, TokenCloseDelim}},
		{"{{!.Active}}", []TokenType{TokenOpenDelim, TokenNot, TokenDot, TokenIdent, TokenCloseDelim}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(tt.input)
			for i, wantType := range tt.want {
				tok, err := l.NextToken()
				if err != nil {
					t.Fatalf("token %d: unexpected error: %v", i, err)
				}

				if tok.Type != wantType {
					t.Errorf("token %d: expected %v, got %v", i, wantType, tok.Type)
				}
			}
		})
	}
}

func TestLexer_StringLiterals(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`{{"hello"}}`, "hello"},
		{`{{'world'}}`, "world"},
		{`{{"Hello, World!"}}`, "Hello, World!"},
		{`{{"with \"quotes\""}}`, `with \"quotes\"`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(tt.input)
			_, _ = l.NextToken() // Skip {{

			tok, err := l.NextToken()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tok.Type != TokenString {
				t.Errorf("expected TokenString, got %v", tok.Type)
			}

			if tok.Value != tt.want {
				t.Errorf("expected %q, got %q", tt.want, tok.Value)
			}
		})
	}
}

func TestLexer_Numbers(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"{{42}}", "42"},
		{"{{123}}", "123"},
		{"{{3.14}}", "3.14"},
		{"{{0.5}}", "0.5"},
		{"{{100.0}}", "100.0"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(tt.input)
			_, _ = l.NextToken() // Skip {{

			tok, err := l.NextToken()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tok.Type != TokenNumber {
				t.Errorf("expected TokenNumber, got %v", tok.Type)
			}

			if tok.Value != tt.want {
				t.Errorf("expected %q, got %q", tt.want, tok.Value)
			}
		})
	}
}

func TestLexer_PipeOperator(t *testing.T) {
	input := "{{.Name | upper | trim}}"
	l := New(input)

	tests := []struct {
		wantType  TokenType
		wantValue string
	}{
		{TokenOpenDelim, "{{"},
		{TokenDot, "."},
		{TokenIdent, "Name"},
		{TokenPipe, "|"},
		{TokenIdent, "upper"},
		{TokenPipe, "|"},
		{TokenIdent, "trim"},
		{TokenCloseDelim, "}}"},
	}

	for i, tt := range tests {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatalf("test %d: unexpected error: %v", i, err)
		}

		if tok.Type != tt.wantType {
			t.Errorf("test %d: expected type %v, got %v", i, tt.wantType, tok.Type)
		}
	}
}

func TestLexer_ArrayAccess(t *testing.T) {
	input := "{{.Items[0]}}"
	l := New(input)

	tests := []TokenType{
		TokenOpenDelim,
		TokenDot,
		TokenIdent, // Items
		TokenLBrack,
		TokenNumber, // 0
		TokenRBrack,
		TokenCloseDelim,
	}

	for i, wantType := range tests {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}

		if tok.Type != wantType {
			t.Errorf("token %d: expected %v, got %v", i, wantType, tok.Type)
		}
	}
}

func TestLexer_ComplexExpression(t *testing.T) {
	input := `{{if .User.Age >= 18 && .User.Active}}`
	l := New(input)

	tokens, err := l.All()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have: {{, if, ., User, ., Age, >=, 18, &&, ., User, ., Active, }}, EOF
	expectedTypes := []TokenType{
		TokenOpenDelim, TokenIf, TokenDot, TokenIdent, TokenDot, TokenIdent,
		TokenGreaterEq, TokenNumber, TokenAnd, TokenDot, TokenIdent, TokenDot,
		TokenIdent, TokenCloseDelim, TokenEOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, want := range expectedTypes {
		if tokens[i].Type != want {
			t.Errorf("token %d: expected %v, got %v", i, want, tokens[i].Type)
		}
	}
}

func TestLexer_LineAndColumn(t *testing.T) {
	input := "Line 1\n{{.Name}}\nLine 3"
	l := New(input)

	// First token: "Line 1\n" starts at line 1
	tok, _ := l.NextToken()
	if tok.Line != 1 {
		t.Errorf("expected line 1, got %d", tok.Line)
	}

	// {{ starts at line 2 (after the newline in previous text)
	tok, _ = l.NextToken()
	if tok.Line != 2 {
		t.Errorf("expected line 2, got %d", tok.Line)
	}

	// Skip to last text token
	for tok.Type != TokenCloseDelim {
		tok, _ = l.NextToken()
	}

	// "\nLine 3" - text starts at position after }}, which is still line 2
	// (the newline is part of the text token value)
	tok, _ = l.NextToken()
	if tok.Line != 2 {
		t.Errorf("expected line 2 for text token starting with newline, got %d", tok.Line)
	}
}

func TestLexer_ErrorUnterminatedString(t *testing.T) {
	input := `{{"unterminated}`
	l := New(input)

	_, _ = l.NextToken() // {{
	_, err := l.NextToken()

	if err == nil {
		t.Error("expected error for unterminated string, got nil")
	}
}

func TestLexer_SpecialLoopVars(t *testing.T) {
	input := "{{@index}} {{@first}} {{@last}}"
	l := New(input)

	// Should tokenize @index, @first, @last as identifiers
	_, _ = l.NextToken() // {{
	tok, _ := l.NextToken()

	if tok.Type != TokenIdent {
		t.Errorf("expected TokenIdent, got %v", tok.Type)
	}

	if tok.Value != "@index" {
		t.Errorf("expected @index, got %q", tok.Value)
	}
}

func BenchmarkLexer_SimpleText(b *testing.B) {
	input := "Hello, World! This is some template text."
	for i := 0; i < b.N; i++ {
		l := New(input)
		_, _ = l.All()
	}
}

func BenchmarkLexer_SimpleVariable(b *testing.B) {
	input := "{{.Name}}"
	for i := 0; i < b.N; i++ {
		l := New(input)
		_, _ = l.All()
	}
}

func BenchmarkLexer_ComplexTemplate(b *testing.B) {
	input := `
		<h1>{{.Title}}</h1>
		{{if .User.Active}}
			<p>Welcome, {{.User.Name | upper}}!</p>
			{{range .Items}}
				<li>{{.}} - {{@index}}</li>
			{{end}}
		{{end}}
	`
	for i := 0; i < b.N; i++ {
		l := New(input)
		_, _ = l.All()
	}
}

func TestToken_String(t *testing.T) {
	tok := Token{Type: TokenIdent, Value: "test", Line: 1, Column: 5}
	s := tok.String()
	if s != "IDENT(\"test\") at 1:5" {
		t.Errorf("unexpected string: %s", s)
	}

	// Test long value truncation
	tok = Token{Type: TokenText, Value: "this is a very long text value that should be truncated", Line: 2, Column: 1}
	s = tok.String()
	if len(s) > 100 {
		t.Errorf("string should be truncated")
	}
}

func TestTokenType_String(t *testing.T) {
	tests := []struct {
		typ  TokenType
		want string
	}{
		{TokenEOF, "EOF"},
		{TokenText, "TEXT"},
		{TokenIdent, "IDENT"},
		{TokenIf, "IF"},
		{TokenPipe, "|"},
		{TokenEqual, "=="},
	}

	for _, tt := range tests {
		got := tt.typ.String()
		if got != tt.want {
			t.Errorf("TokenType.String() = %q, want %q", got, tt.want)
		}
	}
}

func TestLexer_ArithmeticOperators(t *testing.T) {
	input := "{{.A + .B - .C * .D / .E % .F}}"
	l := New(input)

	expectedTypes := []TokenType{
		TokenOpenDelim, TokenDot, TokenIdent, // {{.A
		TokenPlus, TokenDot, TokenIdent, // +.B
		TokenMinus, TokenDot, TokenIdent, // -.C
		TokenMult, TokenDot, TokenIdent, // *.D
		TokenDiv, TokenDot, TokenIdent, // /.E
		TokenMod, TokenDot, TokenIdent, // %.F
		TokenCloseDelim, TokenEOF,
	}

	for i, want := range expectedTypes {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}
		if tok.Type != want {
			t.Errorf("token %d: expected %v, got %v", i, want, tok.Type)
		}
	}
}

func TestLexer_Parentheses(t *testing.T) {
	input := "{{(.A + .B) * .C}}"
	l := New(input)

	expectedTypes := []TokenType{
		TokenOpenDelim, TokenLParen, TokenDot, TokenIdent,
		TokenPlus, TokenDot, TokenIdent, TokenRParen,
		TokenMult, TokenDot, TokenIdent, TokenCloseDelim,
	}

	for i, want := range expectedTypes {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}
		if tok.Type != want {
			t.Errorf("token %d: expected %v, got %v", i, want, tok.Type)
		}
	}
}

func TestLexer_Assignment(t *testing.T) {
	input := "{{name = .User.Name}}"
	l := New(input)

	_, _ = l.NextToken() // {{
	_, _ = l.NextToken() // name

	tok, _ := l.NextToken()
	if tok.Type != TokenAssign {
		t.Errorf("expected TokenAssign, got %v", tok.Type)
	}
}

func TestLexer_MapAccess(t *testing.T) {
	input := `{{.Data["key"]}}`
	l := New(input)

	expectedTypes := []TokenType{
		TokenOpenDelim, TokenDot, TokenIdent, // {{.Data
		TokenLBrack, TokenString, TokenRBrack, // ["key"]
		TokenCloseDelim,
	}

	for i, want := range expectedTypes {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}
		if tok.Type != want {
			t.Errorf("token %d: expected %v, got %v", i, want, tok.Type)
		}
	}
}

func TestLexer_FunctionCall(t *testing.T) {
	input := "{{truncate .Text 100}}"
	l := New(input)

	expectedTypes := []TokenType{
		TokenOpenDelim, TokenIdent, TokenDot, TokenIdent,
		TokenNumber, TokenCloseDelim,
	}

	for i, want := range expectedTypes {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}
		if tok.Type != want {
			t.Errorf("token %d: expected %v, got %v", i, want, tok.Type)
		}
	}
}

func TestLexer_MultipleExpressions(t *testing.T) {
	input := "Hello {{.Name}}, you are {{.Age}} years old"
	l := New(input)

	tokens, err := l.All()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have: TEXT, {{, ., Name, }}, TEXT, {{, ., Age, }}, TEXT, EOF
	if len(tokens) != 12 {
		t.Errorf("expected 12 tokens, got %d", len(tokens))
	}
}

func TestLexer_EmptyTemplate(t *testing.T) {
	l := New("")

	tok, err := l.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tok.Type != TokenEOF {
		t.Errorf("expected TokenEOF, got %v", tok.Type)
	}
}

func TestLexer_OnlyText(t *testing.T) {
	input := "This is just plain text with no expressions"
	l := New(input)

	tok, err := l.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tok.Type != TokenText {
		t.Errorf("expected TokenText, got %v", tok.Type)
	}

	if tok.Value != input {
		t.Errorf("expected %q, got %q", input, tok.Value)
	}
}

func TestLexer_EmptyExpression(t *testing.T) {
	input := "{{}}"
	l := New(input)

	tok1, _ := l.NextToken()
	tok2, _ := l.NextToken()

	if tok1.Type != TokenOpenDelim {
		t.Errorf("expected TokenOpenDelim, got %v", tok1.Type)
	}

	if tok2.Type != TokenCloseDelim {
		t.Errorf("expected TokenCloseDelim, got %v", tok2.Type)
	}
}

func TestLexer_CommaInExpression(t *testing.T) {
	input := "{{fn .Arg1, .Arg2, .Arg3}}"
	l := New(input)

	tokens, err := l.All()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that we have commas
	commaCount := 0
	for _, tok := range tokens {
		if tok.Type == TokenComma {
			commaCount++
		}
	}

	if commaCount != 2 {
		t.Errorf("expected 2 commas, got %d", commaCount)
	}
}
