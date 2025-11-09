package interpreter

import (
	"food-interpreter/generator"
	"food-interpreter/lexer"
	"food-interpreter/nutritionist"
	"food-interpreter/parser"
)

// TODO return the parser for now
func Interpret(diary string) parser.Parser {

	l := lexer.LexString(diary)
	p, nodes := parser.ParseTokens(l.Tokens)

	entries := generator.Generate(nodes)

	nutritionist.AddFoodData(entries)

	return p // TODO return diary entries
}
