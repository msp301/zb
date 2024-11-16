package query

type Lexer struct {
	char     byte
	input    string
	position int
}

func New(input string) *Lexer {
	return &Lexer{
		input: input,
	}
}

func (l *Lexer) readChar() {
	l.char = l.input[l.position]
}

func (l *Lexer) NextToken() Token {
	return Token{
		Type:  TERM,
		Value: "Foo",
	}
}
