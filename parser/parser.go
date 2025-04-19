package parser

import (
	"fmt"
	"food-interpreter/lexer"
)

type Parser struct {
	tokens  []lexer.Token
	current lexer.Token
	index   int
}

// TODO newline seperation? Lexer needs to deal with it if so?
// TODO lexer.Token, is there a better way of doing this?

func ParseTokens(tokens []lexer.Token) Parser {
	parser := Parser{tokens, tokens[0], 0}

	e := parser.parse()
	fmt.Print(e)

	return parser // TODO what exactly should be returned?
}

func (p *Parser) parse() int {
	for p.index < len(p.tokens) {
		p.year() // TODO year isn't the only root
		fmt.Println()
	}
	return 0
}

// Go to the next token
func (p *Parser) nextToken() bool {
	// TODO check for last
	p.index++
	if p.index < len(p.tokens) {
		p.current = p.tokens[p.index]
		return true
	}
	return false
}

// Consume the next token if it's acceptable
func (p *Parser) accept(tokenType int) bool {
	if p.current.TokenType == tokenType {
		fmt.Printf("%v accepted", p.current.TokenType) // TODO print the token type in string not int
		p.nextToken()
		return true
	}
	return false
}

// Check the next token is as expected
func (p *Parser) expect(tokenType int) bool {
	if p.accept(tokenType) {
		return true
	}
	print("Error: unexpected symbol")
	return false
}

func (p *Parser) year() {
	p.accept(lexer.YEAR)
	p.semicolon()
}

//
func (p *Parser) semicolon() {
	p.accept(lexer.SEMICOLON)
}

//
//type Yr interface {
//}

// TODO tree nodes should be structs?
type Year struct {
	Year      lexer.Token
	Semicolon lexer.Token
}

type Semicolon struct {
	Semicolon lexer.Token
}

//type MonthAndDay struct {
//}
