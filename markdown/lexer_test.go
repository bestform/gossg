package markdown_test

import (
	"testing"

	"github.com/bestform/gossg/markdown"
)

func TestInputSameAsOutput(t *testing.T) {
	input := "# Header\n\nHello world"
	lexer := markdown.NewLexer(input)
	tokens := lexer.Tokenize()
	if tokens.String() != input {
		t.Errorf("Expected 'hello world', got '%s'", tokens.String())
	}
}
