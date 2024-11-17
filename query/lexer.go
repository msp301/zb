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
	var tokenValue string

	buffer := bytes.NewBufferString(tokenValue)
	for i := 0; i < len(l.input); i++ {
		switch l.readChar(); l.char {
		default:
			buffer.WriteByte(l.char)
		}
	}

	return Token{
		Type:  TERM,
		Value: buffer.String(),
	}
}
