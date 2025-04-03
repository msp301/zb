package query

import (
	"testing"
)

func TestParse(t *testing.T) {
	lexer := New("test a thing")
	parser := NewParser(lexer)

	// fmt.Println(parser.currentToken.Value)
	// fmt.Println(parser.nextToken.Value)
	// fmt.Println(parser.position)
	//
	// fmt.Println("-------")
	//
	// parser.readToken()
	// fmt.Println(parser.currentToken.Value)
	// fmt.Println(parser.nextToken.Value)
	// fmt.Println(parser.position)
	//
	// fmt.Println("-------")
	//
	// parser.readToken()
	// fmt.Println(parser.currentToken.Value)
	// fmt.Println(parser.nextToken.Value)
	// fmt.Println(parser.nextToken.Type)
	// fmt.Println(parser.position)
	//
	// fmt.Println("-------")
	//
	// parser.readToken()
	// fmt.Println(parser.currentToken.Value)
	// fmt.Println(parser.currentToken.Type)
	// fmt.Println(parser.nextToken.Value)
	// fmt.Println(parser.nextToken.Type)
	// fmt.Println(parser.position)

	parser.Parse()

	t.Fail()
}
