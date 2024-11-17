package query

import (
	"bytes"
)

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
	l.position = l.position + 1
}

func (l *Lexer) NextToken() Token {
	var tokenType TokenType = END
	var tokenValue string

	buffer := bytes.NewBufferString(tokenValue)
READER:
	for l.position < len(l.input) {
		l.readChar()

		if l.char == ' ' {
			break READER
		}

		buffer.WriteByte(l.char)

		switch l.char {
		case '(':
			tokenType = LEFT_BRACKET
			break READER
		default:
			tokenType = TERM
		}
	}

	return Token{
		Type:  tokenType,
		Value: buffer.String(),
	}
}
