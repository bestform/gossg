package markdown

type TokenType int

const (
	T_Hash TokenType = iota
	T_Dash
	T_NewLine
	T_Space
	T_Underscore
	T_Star
	T_Char
	T_EOF
)

type Token struct {
	Value string
	Type  TokenType
}

func (t Token) String() string {
	return t.Value
}

type Tokens []Token

func (t Tokens) String() string {
	var str string
	for _, token := range t {
		str += token.Value
	}
	return str
}

type RuneParser func(lexer *Lexer) *Token

func singleCharParser(char rune, tokenType TokenType) RuneParser {
	return func(lexer *Lexer) *Token {
		if lexer.data[lexer.pos] == char {
			lexer.pos++
			return &Token{Value: string(char), Type: tokenType}
		}
		return nil
	}
}

func anyCharParser(lexer *Lexer) *Token {
	char := lexer.data[lexer.pos]
	lexer.pos++
	return &Token{Value: string(char), Type: T_Char}
}

type Lexer struct {
	parsers []RuneParser
	data    []rune
	pos     int
}

func NewLexer(data string) *Lexer {
	return &Lexer{
		data: []rune(data),
		parsers: []RuneParser{
			singleCharParser('#', T_Hash),
			singleCharParser(' ', T_Space),
			singleCharParser('-', T_Dash),
			singleCharParser('\n', T_NewLine),
			singleCharParser('_', T_Underscore),
			singleCharParser('*', T_Star),
			anyCharParser,
		},
	}
}

func (l *Lexer) Tokenize() Tokens {
	var tokens Tokens
	for {
		token := l.NextToken()
		tokens = append(tokens, *token)
		if token.Type == T_EOF {
			break
		}
	}

	return tokens
}

func (l *Lexer) NextToken() *Token {
	if l.pos >= len(l.data) {
		return &Token{Type: T_EOF}
	}

	for _, parser := range l.parsers {
		if token := parser(l); token != nil {
			return token
		}
	}

	return nil
}
