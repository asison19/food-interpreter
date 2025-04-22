package main

import (
	"food-interpreter/lexer"
	"food-interpreter/parser"
)

func main() {
	//lex := lexer.LexFile("test/data/year.txt")
	lex := lexer.LexFile("test/data/test.txt")

	parser.ParseTokens(lex.Tokens)
}
