package query

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
