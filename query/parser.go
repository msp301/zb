package query

import (
	"fmt"
	"strings"

	"github.com/msp301/graph"
)

type Parser struct {
	lexer        *Lexer
	currentToken Token
	nextToken    Token
	position     int
}

func NewParser(lexer *Lexer) *Parser {
	p := Parser{
		lexer: lexer,
	}
	p.readToken()
	p.readToken()
	p.position = 0
	return &p
}

func (p *Parser) readToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
	p.position = p.position + 1
}

func (p *Parser) Parse() *graph.Graph {
	ast := graph.Directed()

	for p.currentToken.Type != END {
		ast.Add(uint64(p.position+1), p.currentToken.String(), p.currentToken)
		p.readToken()
	}

	ast.Walk(func(vertex graph.Vertex, depth int) bool {
		indent := strings.Repeat("\t", depth)
		fmt.Printf("%s%s\n", indent, vertex.Label)
		return true
	}, -1)

	return ast
}
