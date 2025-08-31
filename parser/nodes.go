package parser

import (
	"food-interpreter/lexer"
)

type year struct {
	year      lexer.Token
	semicolon semicolon
}

type semicolon struct {
	semicolon lexer.Token
}
