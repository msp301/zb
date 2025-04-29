package query

import "fmt"

const (
	TERM = iota
	AND
	OR
	NOT
	LEFT_BRACKET
	RIGHT_BRACKET
	SINGLE_QUOTE
	DOUBLE_QUOTE
	END
)

type TokenType int

type Token struct {
	Type  TokenType
	Value string
}

func (t *Token) String() string {
	return tokenString(t.Type)
}

func tokenString(tokenType TokenType) string {
	switch tokenType {
	case TERM:
		return "TERM"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case LEFT_BRACKET:
		return "("
	case RIGHT_BRACKET:
		return ")"
	case SINGLE_QUOTE:
		return "'"
	case DOUBLE_QUOTE:
		return "\""
	case END:
		return "END"
	default:
		panic(fmt.Sprintf("Unknown token: %d", tokenType))
	}
}
