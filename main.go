package main

import (
	"bufio"
	"fmt"
	"food-interpreter/lex"
	"os"
)

func main() {

	args := os.Args[1:]

	// TODO command line args
	scanFile(args[0])
}

// TODO error handling
func scanFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lex.ScanTokens(reader)
}
