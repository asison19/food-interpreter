package interpreter

import (
	"food-interpreter/generator"
	"food-interpreter/lexer"
	"food-interpreter/parser"
)

// TODO return the parser for now
func Interpret(diary string) parser.Parser {

	l := lexer.LexString(diary)
	p, nodes := parser.ParseTokens(l.Tokens)

	generator.Generate(nodes)

	return p // TODO return diary entries
}
