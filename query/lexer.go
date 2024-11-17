package query

import (
	"bytes"
	"slices"
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
			if len(buffer.String()) == 0 {
				continue
			}

			switch buffer.String() {
			case "and":
				tokenType = AND
				break
			case "or":
				tokenType = OR
				break
			case "not":
				tokenType = NOT
				break
			}

			break READER
		}

		terminatingChars := []byte{'(', ')', '\'', '"'}
		if len(buffer.String()) > 0 && slices.Contains(terminatingChars, l.char) {
			l.position = l.position - 1
			break READER
		}

		buffer.WriteByte(l.char)

		switch l.char {
		case '(':
			tokenType = LEFT_BRACKET
			break READER
		case ')':
			tokenType = RIGHT_BRACKET
			break READER
		case '\'':
			tokenType = SINGLE_QUOTE
			break READER
		case '"':
			tokenType = DOUBLE_QUOTE
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
