package parser

import (
	"fmt"
	"food-interpreter/lexer"
)

type Parser struct {
	Tokens  []lexer.Token
	current lexer.Token
	index   int
	//errorCount int // TODO, only do up to 3-5 errors, probably best to go straight to semicolon, or the next line
}

// TODO newline seperation? Lexer needs to deal with it if so?
// TODO lexer.Token, is there a better way of doing this?

func ParseTokens(tokens []lexer.Token) Parser {
	parser := Parser{tokens, tokens[0], 0}

	parser.parse()

	return parser // TODO what exactly should be returned? Create AST nodes and return them on each function, then return the root (need to make a complete root?) here?
}

func (p *Parser) parse() int {
	for p.index < len(p.Tokens) {
		token, _ := p.check()
		switch token.Type {
		case lexer.YEAR:
			p.year()
		case lexer.MONTHANDDAY:
			p.monthAndDay()
		default:
			fmt.Printf("Year or MonthAndDay expected, got %v instead", token.Type)
			p.nextToken() // We're allowing a continue
		}
		fmt.Println()
	}
	return 0
}

// Go to the next token
func (p *Parser) nextToken() bool {
	p.index++
	if p.index < len(p.Tokens) {
		p.current = p.Tokens[p.index]
		return true
	}
	return false
}

// Accept the current token if it's the same as the passed in token.
func (p *Parser) accept(tokenType lexer.TokenType) bool {
	if p.current.Type == tokenType {
		fmt.Printf("%v accepted", p.current.Type)
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
	fmt.Printf("Error: unexpected symbol %v", tokenType)
	return false
}

func (p *Parser) year() {
	if p.accept(lexer.YEAR) {
		p.semicolon() // TODO semicolon should be optional. At terminals there should be the error?
		return
	}
	fmt.Printf("Year expected, got %v instead", p.Tokens[p.index])
	p.nextToken() // We're allowing a continue
}

func (p *Parser) check() (lexer.Token, bool) {
	if len(p.Tokens) <= p.index {
		return lexer.Token{}, true
	}
	return p.Tokens[p.index], false
}

// TODO index out of range for 01/23 last in token slice
func (p *Parser) monthAndDay() {
	p.expect(lexer.MONTHANDDAY)
	p.time()
}

func (p *Parser) time() {
	p.expect(lexer.TIME)

	token, diary_err := p.check()
	if diary_err {
		fmt.Printf("Received time token and expected food, repeater, or sleep next.")
		return
	}
	switch token.Type {
	case lexer.FOOD:
		p.food()
	case lexer.REPEATER:
		p.repeater()
	case lexer.SLEEP:
		p.sleep()
	default:
		fmt.Printf("Food, repeater, or sleep expected, got %v instead", token.Type)
		p.nextToken()
	}

}
func (p *Parser) food() {
	p.expect(lexer.FOOD)
	token, diary_err := p.check()
	if diary_err {
		fmt.Printf("Received food token and expected comma or semicolon next.")
		return
	}
	switch token.Type {
	case lexer.COMMA: // TODO rename comma nonterminal
		p.comma()
	case lexer.SEMICOLON: // TODO turn semicolon to a terminal?
		p.semicolon()
	default:
		fmt.Printf("Comma or semicolon expected, got %v instead", token.Type)
		p.nextToken()
	}
}
func (p *Parser) repeater() {
	p.expect(lexer.REPEATER)
	token, diary_err := p.check()
	if diary_err {
		fmt.Printf("Received repeater token and expected comma or semicolon next.")
		return
	}
	switch token.Type {
	case lexer.COMMA:
		p.comma()
	case lexer.SEMICOLON:
		p.semicolon()
	default:
		fmt.Printf("Comma or semicolon expected, got %v instead", token.Type)
		p.nextToken()
	}
}
func (p *Parser) sleep() {
	p.expect(lexer.SLEEP)
	token, diary_err := p.check()
	if diary_err {
		fmt.Printf("Received sleep token and expected comma or semicolon next.")
		return
	}
	switch token.Type {
	case lexer.COMMA:
		p.comma()
	case lexer.SEMICOLON:
		p.semicolon()
	default:
		fmt.Printf("Comma or semicolon expected, got %v instead", token.Type)
		p.nextToken()
	}
}

// TODO Update the language. After a comma could be a repeater, or sleep.
func (p *Parser) comma() {
	p.expect(lexer.COMMA)
	p.food()
}

// TODO do I really need this to be a nonterminal?
// TODO make it so that it's optional for the end of a line for MonthAndDay root or at least just Year.
func (p *Parser) semicolon() {
	p.accept(lexer.SEMICOLON)
}
