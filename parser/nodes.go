package parser

import (
	"food-interpreter/lexer"
)

type Node interface{}

type Year struct {
	year      lexer.Token
	semicolon Node
}

type Semicolon struct {
	semicolon lexer.Token
}
