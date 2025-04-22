package main

import (
	"fmt"
	"food-interpreter/lexer"
	"food-interpreter/parser"
)

func main() {
	//lex := lexer.LexFile("test/data/year.txt")
	lex := lexer.LexFile("test/data/test.txt")

	par := parser.ParseTokens(lex.Tokens)
}
