package parser

import (
	"testing"

	"github.com/toutaio/toutago-fith-renderer/lexer"
)

func TestParserError(t *testing.T) {
	// Test parser error handling
	l := lexer.New("{{.invalid syntax")
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Error("Expected parse error for invalid syntax")
	}
}
