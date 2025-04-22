package parser

import (
	"fmt"
	"food-interpreter/lexer"
)

type Parser struct {
	tokens  []lexer.Token
	current lexer.Token
	index   int
	//errorCount int // TODO, only do up to 3-5 errors, probably best to go straight to semicolon, or the next line
}

// TODO newline seperation? Lexer needs to deal with it if so?
// TODO lexer.Token, is there a better way of doing this?

func ParseTokens(tokens []lexer.Token) Parser {
	parser := Parser{tokens, tokens[0], 0}

	parser.parse()

	return parser // TODO what exactly should be returned?
}

func (p *Parser) parse() int {
	for p.index < len(p.tokens) {
		switch p.check().Type {
		case lexer.YEAR:
			p.year() // TODO year isn't the only root. p.check() it?
		case lexer.MONTHANDDAY:
			p.monthAndDay() // TODO year isn't the only root. p.check() it?
		default:
			fmt.Printf("Year or MonthAndDay expected, got %v instead", p.check())
			p.nextToken() // We're allowing a continue
		}
		fmt.Println()
	}
	return 0
}

// Go to the next token
// TODO gracefully handling if it's the last token
func (p *Parser) nextToken() bool {
	// TODO check for last
	p.index++
	if p.index < len(p.tokens) {
		p.current = p.tokens[p.index]
		return true
	}
	return false
}

// Accept the current token if it's the same as the passed in token.
func (p *Parser) accept(tokenType lexer.TokenType) bool {
	if p.current.Type == tokenType {
		fmt.Printf("%v accepted", p.current.Type) // TODO print the token type in string not int
		p.nextToken()
		return true
	}
	return false
}

// The passed in token is the expected current (unconsumed) token
// if not, that's a syntax error.
func (p *Parser) expect(tokenType lexer.TokenType) bool {
	if p.accept(tokenType) {
		return true
	}
	fmt.Printf("Error: unexpected symbol %v", p.check())
	return false
}

func (p *Parser) year() {
	if p.accept(lexer.YEAR) {
		p.semicolon() // TODO semicolon should be optional. At terminals there should be the error?
		return
	}
	fmt.Printf("Year expected, got %v instead", p.check())
	p.nextToken() // We're allowing a continue
}

func (p *Parser) check() lexer.Token {
	return p.tokens[p.index]
}

// TODO index out of range for 01/23 last in token slice
func (p *Parser) monthAndDay() {
	p.expect(lexer.MONTHANDDAY)
	p.time()
}

func (p *Parser) time() {
	p.expect(lexer.TIME)

	// TODO is this the best way of going about this?
	switch p.check().Type {
	case lexer.FOOD:
		p.food()
	case lexer.REPEATER:
		p.repeater()
	case lexer.SLEEP:
		p.sleep()
	default:
		fmt.Printf("Food, repeater, or sleep expected, got %v instead", p.check())
		p.nextToken()
	}

}
func (p *Parser) food() {
	p.expect(lexer.FOOD)
	switch p.check().Type {
	case lexer.COMMA: // TODO rename comma nonterminal
		p.comma()
	case lexer.SEMICOLON: // TODO turn semicolon to a terminal?
		p.semicolon()
	default:
		fmt.Printf("Comma, or semicolon expected, got %v instead", p.check())
		p.nextToken()
	}
}
func (p *Parser) repeater() {
	p.expect(lexer.REPEATER)
	switch p.check().Type {
	case lexer.COMMA:
		p.comma()
	case lexer.SEMICOLON:
		p.semicolon()
	default:
		fmt.Printf("Comma, or semicolon expected, got %v instead", p.check())
		p.nextToken()
	}
}
func (p *Parser) sleep() {
	p.expect(lexer.SLEEP)
	switch p.check().Type {
	case lexer.COMMA:
		p.comma()
	case lexer.SEMICOLON:
		p.semicolon()
	default:
		fmt.Printf("Comma, or semicolon expected, got %v instead", p.check())
		p.nextToken()
	}
}

// TODO Update the language. After a comma could be a repeater, or sleep.
func (p *Parser) comma() {
	p.expect(lexer.COMMA)
	p.food()
}

//
func (p *Parser) semicolon() {
	p.accept(lexer.SEMICOLON)
}
