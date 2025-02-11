package markdown_test

import (
	"testing"

	"github.com/bestform/gossg/markdown"
)

func TestSimpleInput(t *testing.T) {
	input := "## Hello World!\n---\nfoo - bar\nbar"
	lexer := markdown.NewLexer(input)
	tokens := lexer.Tokenize()
	nodes := markdown.NewParser(tokens).Parse()

	expected := "<h2>Hello World!</h2><hr>foo - bar<br>bar"
	if nodes.Render() != expected {
		t.Errorf("Expected %s, got %s", expected, nodes.Render())
	}
}
