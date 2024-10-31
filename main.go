package main

import (
	"bufio"
	"fmt"
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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
