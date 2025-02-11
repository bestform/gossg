package markdown

import (
	"fmt"
)

type NodeType int

const (
	N_Heading NodeType = iota
	N_Text
	N_NewLine
	N_Line
	N_Space
)

type Node struct {
	Value string
	Type  NodeType
	Level int
}

func (n Node) Render() string {
	switch n.Type {
	case N_Heading:
		return fmt.Sprintf("<h%d>%s</h%d>", n.Level, n.Value, n.Level)
	case N_Line:
		return "<hr>"
	case N_NewLine:
		return "<br>"
	case N_Text:
		return n.Value
	case N_Space:
		return " "
	default:
		panic(fmt.Sprintf("unexpected markdown.NodeType: %#v", n.Type))
	}
}

type Nodes []Node

func (n Nodes) Render() string {
	result := ""
	for _, node := range n {
		result += node.Render()
	}

	return result
}

type Parser struct {
	tokens      Tokens
	nodeParsers []nodeParser
	pos         int
}

func NewParser(tokens Tokens) *Parser {
	return &Parser{
		tokens: tokens,
		nodeParsers: []nodeParser{
			headingParser,
			lineParser,
			newLineParser,
			spaceParser,
			textParser,
		},
	}
}

func (p *Parser) Parse() Nodes {
	nodes := Nodes{}
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type != T_EOF {
		foundNode := false
		start := p.pos
		for _, parser := range p.nodeParsers {
			if node := parser(p); node != nil {
				foundNode = true
				nodes = append(nodes, *node)
				break
			}
		}
		if !foundNode {
			panic(fmt.Sprintf("unexpected token: %#v", p.tokens[p.pos]))
		}
		if start == p.pos {
			panic(fmt.Sprintf("no progress: %#v", p.tokens[p.pos]))
		}
	}

	return nodes
}

type nodeParser func(p *Parser) *Node

func newLineParser(p *Parser) *Node {
	if p.tokens[p.pos].Type != T_NewLine {
		return nil
	}
	p.pos++

	return &Node{Type: N_NewLine}
}

func spaceParser(p *Parser) *Node {
	if p.tokens[p.pos].Type != T_Space {
		return nil
	}
	p.pos++

	return &Node{Type: N_Space}
}

func lineParser(p *Parser) *Node {
	dashes := 0
	for p.tokens[p.pos].Type == T_Dash {
		dashes++
		p.pos++
	}
	if dashes < 3 || p.tokens[p.pos].Type != T_NewLine {
		p.pos -= dashes
		return nil
	}
	// swallow the newline
	p.pos++
	return &Node{Type: N_Line}
}

func textParser(p *Parser) *Node {
	text := ""
	t := p.tokens[p.pos].Type
	for t != T_NewLine && t != T_EOF {
		text += p.tokens[p.pos].Value
		p.pos++
		t = p.tokens[p.pos].Type
	}

	return &Node{Value: text, Type: N_Text}
}

func headingParser(p *Parser) *Node {
	level := 0
	for p.tokens[p.pos].Type == T_Hash {
		level++
		p.pos++
	}
	if level == 0 {
		return nil
	}
	if p.tokens[p.pos].Type != T_Space {
		p.pos -= level
		return nil
	}
	// swallow the space
	p.pos++
	text := ""
	for p.tokens[p.pos].Type != T_NewLine && p.tokens[p.pos].Type != T_EOF {
		text += p.tokens[p.pos].Value
		p.pos++
	}

	if p.tokens[p.pos].Type == T_NewLine {
		p.pos++
	}

	return &Node{Value: text, Type: N_Heading, Level: level}
}
