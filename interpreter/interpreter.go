package interpreter

import (
	"food-interpreter/lexer"
	"food-interpreter/parser"
)

// TODO return the parser for now
func Interpret(diary string) parser.Parser {

	l := lexer.LexString(diary)
	p := parser.ParseTokens(l.Tokens)

	return p
}
