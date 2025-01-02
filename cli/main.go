package main

import (
	"food-interpreter/lexer"
	"os"
)

func main() {
	args := os.Args[1:]
	lexer.LexFile(args[0])
}
