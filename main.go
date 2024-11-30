package main

import (
	"food-interpreter/lexer"
	"os"
)

func main() {

	args := os.Args[1:]

	// TODO command line args
	lexer.LexFile(args[0])
}
