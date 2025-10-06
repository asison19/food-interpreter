package interpreter

import (
	"fmt"
	"food-interpreter/generator"
	"food-interpreter/lexer"
	"food-interpreter/parser"
)

// TODO return the parser for now
func Interpret(diary string) parser.Parser {

	l := lexer.LexString(diary)
	p, nodes := parser.ParseTokens(l.Tokens)

	entries := generator.Generate(nodes)

	fmt.Println(entries)

	return p // TODO return diary entries
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
