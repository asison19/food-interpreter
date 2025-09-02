package parser

import (
	"food-interpreter/lexer"
)

type Node interface {
	accept(lexer.Token) Node
	//expect() bool
}

type Year struct {
	year      lexer.Token
	semicolon Node
}

type Semicolon struct {
	semicolon lexer.Token
}

//	func accept(n Node) Node {
//		return n
//	}
//func expect(current lexer.TokenType, expected lexer.TokenType) bool {
//	return current == expected
//}

func (y Year) accept(current lexer.Token) Node {
	if current.Type == lexer.YEAR {
		return Year{current, Semicolon{}} // TODO use legit Semicolon
	}
	return Year{}
}
func (s Semicolon) accept(current lexer.Token) Node {
	if current.Type == lexer.SEMICOLON {
		return Semicolon{current}
	}
	return Semicolon{}
}

//func (y *Year) accept() {
//
//}
