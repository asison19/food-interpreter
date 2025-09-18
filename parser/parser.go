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

func ParseTokens(tokens []lexer.Token) (Parser, []Node) {
	parser := Parser{tokens, tokens[0], 0}

	nodes := parser.parse()

	return parser, nodes
}

// TODO return errors as well
func (p *Parser) parse() []Node {
	var nodes []Node
	for p.index < len(p.Tokens) {
		token, _ := p.check()
		switch token.Type {
		case lexer.YEAR:
			nodes = append(nodes, p.year())
		case lexer.MONTHANDDAY:
			nodes = append(nodes, p.monthAndDay())
		default:
			fmt.Printf("Year or MonthAndDay expected, got %v instead", token.Type)
			p.nextToken() // We're allowing a continue
		}
		fmt.Println()
	}
	return nodes
}

// Go to the next token
// TODO make this better
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
	fmt.Printf("Error: expected symbol %v", tokenType)
	return false
}

func (p *Parser) year() Year {
	y := p.current // TODO there has to be a better way
	p.expect(lexer.YEAR)
	return Year{y, p.semicolon()}
}

// Returns the token and true if there exists another token
// Returns an empty token and false if not
func (p *Parser) check() (lexer.Token, bool) {
	if len(p.Tokens) <= p.index {
		return lexer.Token{}, true
	}
	return p.Tokens[p.index], false
}

func (p *Parser) monthAndDay() MonthAndDay {
	mad := p.current
	p.expect(lexer.MONTHANDDAY)
	return MonthAndDay{mad, p.time()}
}

func (p *Parser) time() Time {
	t := p.current
	if !p.expect(lexer.TIME) {
		fmt.Printf("Error: expected time, got %v instead.", p.current)
	}

	token, diary_err := p.check()
	if diary_err {
		fmt.Printf("Received time token and expected food, repeater, or sleep next.")
		return Time{}
	}
	// TODO simplify by checking for any of the three, go to the next symbol (consume), and then expect comma or semicolon again?
	switch token.Type {
	case lexer.FOOD:
		return Time{t, p.food()}
	case lexer.REPEATER:
		return Time{t, p.repeater()}
	case lexer.SLEEP:
		return Time{t, p.sleep()}
	default:
		fmt.Printf("Food, repeater, or sleep expected, got %v instead", token.Type)
		p.nextToken()
	}
	return Time{}

}
func (p *Parser) food() Food {
	f := p.current
	p.expect(lexer.FOOD)
	token, diary_err := p.check()
	if diary_err {
		fmt.Printf("Received food token and expected comma or semicolon next.")
		return Food{}
	}
	switch token.Type {
	case lexer.COMMA:
		return Food{f, p.comma()}
	case lexer.SEMICOLON:
		return Food{f, p.semicolon()}
	default:
		fmt.Printf("Comma or semicolon expected, got %v instead", token.Type)
		p.nextToken()
		return Food{}
	}
}
func (p *Parser) repeater() Repeater {
	r := p.current
	p.expect(lexer.REPEATER)
	token, diary_err := p.check()
	if diary_err {
		fmt.Printf("Received repeater token and expected comma or semicolon next.")
		return Repeater{}
	}
	switch token.Type {
	case lexer.COMMA:
		return Repeater{r, p.comma()}
	case lexer.SEMICOLON:
		return Repeater{r, p.semicolon()}
	default:
		fmt.Printf("Comma or semicolon expected, got %v instead", token.Type)
		p.nextToken()
		return Repeater{}
	}
}
func (p *Parser) sleep() Sleep {
	s := p.current
	p.expect(lexer.SLEEP)
	token, diary_err := p.check()
	if diary_err {
		fmt.Printf("Received sleep token and expected comma or semicolon next.")
		return Sleep{}
	}
	switch token.Type {
	case lexer.COMMA:
		return Sleep{s, p.comma()}
	case lexer.SEMICOLON:
		return Sleep{s, p.semicolon()}
	default:
		fmt.Printf("Comma or semicolon expected, got %v instead", token.Type)
		p.nextToken()
		return Sleep{}
	}
}

func (p *Parser) comma() Comma {
	c := p.current
	p.expect(lexer.COMMA)
	token, diary_err := p.check()
	if diary_err {
		fmt.Printf("Received comma token and expected food, repeater or sleep next.")
		return Comma{}
	}
	switch token.Type {
	case lexer.FOOD:
		return Comma{c, p.food()}
	//case lexer.REPEATER: // TODO after comma has to come a food, revisit this?
	//	return Comma{c, p.repeater()}
	//case lexer.SLEEP:
	//	return Comma{c, p.sleep()}
	default:
		fmt.Printf("Food, repeater, or sleep expected, got %v instead", token.Type)
		p.nextToken()
		return Comma{}
	}
}

func (p *Parser) semicolon() Semicolon {
	s := p.current
	p.expect(lexer.SEMICOLON)
	token, _ := p.check()
	if token.Type == lexer.TIME {
		return Semicolon{s, p.time()}
	}
	return Semicolon{s, Time{}}
}
