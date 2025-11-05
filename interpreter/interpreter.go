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

	foods := generator.Generate(nodes)

	nutritionist.AddFoodData(foods)

	return p // TODO return diary entries
}
